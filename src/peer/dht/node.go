package dht

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"
	"tp/common"
	"tp/peer/dht/bucket_table"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/dht/key_value_table"
	"tp/peer/dht/tiered_contact_storage"

	"tp/peer/helpers"
	"tp/peer/helpers/communication/rpc_ops"
	"tp/peer/helpers/file_manager"
	"tp/peer/helpers/file_manager/blocks"
)

const MSG_ERROR_OWN_REQUEST = "it is my own request"
const MSG_MUST_DISCARD_CONTACT = "Contact request should be discarded: %v"
const MSG_ERROR_FILE_NOT_FOUND = "the file could not be found: %v"
const MSG_FILE_FOUND = "the file has been found: %v"
const MSG_FILE_DOWLOADED = "the file has been fully downloaded: %v"
const MSG_CONTACTS_FOUND_FOR_KEY = "%v contacts found for key %v"
const MAX_CHAN_PENDING_CONTACTS = 100

// Representa un nodo de una Distributed Hash Table
type Node struct {
	Config                helpers.PeerConfig
	BucketTab             bucket_table.BucketTable
	KeyValueTab           key_value_table.KeyValueTable
	SndStore              rpc_ops.StoreOp
	SndShareContactsRecip rpc_ops.SndShareContactsRecipOp
	SndPing               rpc_ops.PingOp
	SndFindBlock          rpc_ops.FindBlockOp
	PendingContactsToAdd  chan contacts_queue.Contact
}

// Retorna una nueva instancia de nodo lista para ser utilizada
func NewNode(
	config helpers.PeerConfig,
	sndPing rpc_ops.PingOp,
	sndStore rpc_ops.StoreOp,
	sndShareContactsRecip rpc_ops.SndShareContactsRecipOp,
	sndFindBlock rpc_ops.FindBlockOp) *Node {
	node := &Node{
		Config:                config,
		BucketTab:             *bucket_table.NewBucketTable(config, sndPing),
		KeyValueTab:           *key_value_table.NewKeyValueTable(),
		SndStore:              sndStore,
		SndShareContactsRecip: sndShareContactsRecip,
		SndPing:               sndPing,
		SndFindBlock:          sndFindBlock,
		PendingContactsToAdd:  make(chan contacts_queue.Contact, MAX_CHAN_PENDING_CONTACTS),
	}
	return node
}

// Se eliminan recursos asociados
func (node *Node) DisposeNode() {
	close(node.PendingContactsToAdd)
}

// Retorna verdadero si la instancia es el bootstrap node
func (node *Node) IsBootstrapNode() bool {
	return bytes.Equal(node.Config.Id, helpers.BootstrapNodeID)
}

// Chequea si es un contacto a si mismo e intenta agregarlo
func (node *Node) AddContactPreventingLoop(sourceContact contacts_queue.Contact) bool {
	// Prevenir bucle
	if node.discardContact(sourceContact) {
		common.Log.Debugf(fmt.Sprintf(MSG_MUST_DISCARD_CONTACT, sourceContact.ToString()))
		return false
	}
	// Trata de agregar el contacto
	node.BucketTab.AddContact(sourceContact)
	return true
}

// Representa la recepción de un ping el cuál consiste en intentar agregar el contacto a la tabla de
// contactos
func (node *Node) RcvPing(sourceContact contacts_queue.Contact) bool {
	return node.AddContactPreventingLoop(sourceContact)
}

// Obtiene los contactos locales recomendados para la fuente, agrega los contactos compartidos por la fuente y
// retorna los contactos recomendados para la fuente
func (node *Node) RcvShCtsRecip(sourceContact contacts_queue.Contact, sourceContactList []contacts_queue.Contact) []contacts_queue.Contact {
	// agregar contacto origen
	node.AddContactPreventingLoop(sourceContact)
	// obtener contactos recomendados
	selfContacts := node.BucketTab.GetRecommendedContactsForId(sourceContact.ID)
	// agregar contactos que compartió la fuente
	node.BucketTab.AddContacts(selfContacts)
	return selfContacts
}

// Envía los contactos propios al bootstrap node esperando que el mismo retorne los contactos recomendados
// para la clave del presente nodo
func (node *Node) SndShCtsToBootstrap() {
	if !node.IsBootstrapNode() {
		contactBoostrapNode := contacts_queue.NewContact(helpers.BootstrapNodeID, helpers.BootstrapNodeUrl)
		// agregar bootstrap node a contactos
		if node.SndShCts(*contactBoostrapNode) == nil {
			node.AddContactPreventingLoop(*contactBoostrapNode)
		}
	}
}

// Envía los contactos propios al contacto node esperando que el mismo retorne los contactos recomendados
// para la clave del presente nodo
func (node *Node) SndShCts(destContact contacts_queue.Contact) error {
	// ¿Esta vivo el nodo?
	err := node.SndPing(node.Config, destContact)
	if err != nil {
		return err
	}
	// obtener contactos recomendados
	selfContacts := node.BucketTab.GetRecommendedContactsForId(destContact.ID)
	// enviar contactos a contacto desitno
	destRcvContacts, err := node.SndShareContactsRecip(node.Config, destContact, selfContacts)
	if err != nil {
		return err
	}
	// agregar contactos recibidos
	//node.AddContactsCheckingState(destRcvContacts)
	node.addContactsDefferedPing(destRcvContacts)
	//node.BucketTab.AddContacts(destRcvContacts)
	return nil
}

// Retorna los contactos de los nodos más cercanos a un targetId. Además hace el intento de
// agregar el contacto solicitante a la bucket_table
func (node *Node) RcvFindNode(sourceContact contacts_queue.Contact, targetId []byte) []contacts_queue.Contact {
	// Prevenir bucle
	if node.discardContact(sourceContact) {
		common.Log.Debugf(fmt.Sprintf(MSG_MUST_DISCARD_CONTACT, sourceContact.ToString()))
		return []contacts_queue.Contact{}
	}
	// Agregar contacto a la bucket_table
	node.AddContactPreventingLoop(sourceContact)
	// Buscar los contactos
	return node.BucketTab.GetContactsForId(targetId)
}

// Si la target key se encuentra en el nodo retorna el archivo asociado a la misma, caso contrario retorna
// un error y la lista de los contactos más cercanos a la misma. Además hace el intento de
// agregar el contacto solicitante a la bucket_table
func (node *Node) RcvFindBlock(sourceContact contacts_queue.Contact, targetKey []byte) (string, []byte, []contacts_queue.Contact, error) {
	// Prevenir bucle
	if node.discardContact(sourceContact) {
		common.Log.Debugf(fmt.Sprintf(MSG_MUST_DISCARD_CONTACT, sourceContact.ToString()))
		return key_value_table.EMPTY_VALUE, []byte{}, []contacts_queue.Contact{}, errors.New(MSG_ERROR_OWN_REQUEST)
	}
	// Agregar contacto a la bucket_table
	node.AddContactPreventingLoop(sourceContact)
	// Búsqueda de archivo
	fileName, data, err := node.KeyValueTab.Get(targetKey)
	if err == nil {
		return fileName, data, []contacts_queue.Contact{}, nil
	}
	contactsToReturn := node.BucketTab.GetContactsForId(targetKey)
	common.Log.Debugf(MSG_CONTACTS_FOUND_FOR_KEY, len(contactsToReturn), helpers.KeyToHexString(targetKey))
	return key_value_table.EMPTY_VALUE, []byte{}, contactsToReturn, nil
}

// Almacena la clave valor localmente y envía el menseja de store a los contactos más cercanos a la tabla.
// En caso de que la clave ya existía localmente retorna error. Por otro lado intenta agregar el contacto
// fuente en la tabla de contactos
func (node *Node) RcvStore(sourceContact contacts_queue.Contact, key []byte, fileName string, data []byte) error {
	// Prevenir bucle
	if node.discardContact(sourceContact) {
		common.Log.Debugf(fmt.Sprintf(MSG_MUST_DISCARD_CONTACT, sourceContact.ToString()))
		return errors.New(MSG_ERROR_OWN_REQUEST)
	}
	// Agregar contacto a la bucket_table
	node.AddContactPreventingLoop(sourceContact)
	// Almacenar localmente
	return node.doStoreBlock(key, fileName, data)
}

// Intenta guardar bloque localmente y añadir key a la tabla
func (node *Node) doStoreBlock(key []byte, fileName string, data []byte) error {
	// Almacenar localmente
	err := node.KeyValueTab.Add(key, fileName, data)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	// Buscar contactos cercanos a la clave
	contacts := node.BucketTab.GetContactsForId(key)
	for index := range contacts {
		node.SndStore(node.Config, contacts[index], key, fileName, data)
	}
	return nil
}

// Retorna verdadero si la url propia y la del contacto coinciden
func (node *Node) discardContact(contact contacts_queue.Contact) bool {
	return node.Config.Url == contact.Url
}

// Retorna los contactos para un id dado
func (node *Node) getContactsForId(id []byte) []contacts_queue.Contact {
	return node.BucketTab.GetContactsForId(id)
}

// Agrega un archivo del espacio local al ipfs dado por los nodos de la red de contactos
func (node *Node) AddFile(fileName string) error {
	return file_manager.AddFile(fileName, node.createSndBlockNeighbors())
}

// Busca el archivo localmente y en la red de nodos
func (node *Node) GetFile(fileName string) error {
	// Primer bloque
	blockName := blocks.GenerateBlockName(fileName, 0)
	key := helpers.GetKey(blockName)
	// Obtener archivo completo
	endFile := false
	for !endFile {
		// Buscar bloque localmente
		if end, nextBlockKey, err := node.findBlockLocally(key); err == nil {
			key = nextBlockKey
			endFile = end
			continue
		}
		// Obtengo contactos locales
		localContacts := node.getContactsForId(key)
		// Si no hay contactos locales se retorna error
		if len(localContacts) == 0 {
			msg := fmt.Sprintf(MSG_ERROR_FILE_NOT_FOUND, fileName)
			common.Log.Errorf(msg)
			return errors.New(msg)
		}
		// creo el storage de contactos y agrego los locales
		contactStorage := tiered_contact_storage.NewTieredContactStorage(key)
		contactStorage.PushContacts(localContacts)
		// Obtener bloque
		endBlock := false
		for !endBlock {
			var errorToReturn error = nil
			numThreads := 20
			resChan := make(chan processNextContactReturn, numThreads)
			defer close(resChan)
			wg := new(sync.WaitGroup)
			// Procesar contacto
			for id := range numThreads {
				wg.Add(1)
				go func() {
					processNextContact(node, key, fileName, contactStorage, resChan, id)
					wg.Done()
				}()
			}
			wg.Wait()
			for range numThreads {
				resProc := <-resChan
				if !resProc.blockHasBeenFound {
					if resProc.err != nil {
						errorToReturn = resProc.err
					}
					continue
				}
				endFile = resProc.endFileFound
				key = resProc.nextBlockFound
				endBlock = true
				break
			}
			if !endBlock && errorToReturn != nil {
				return errorToReturn
			}
		}
		if endFile {
			common.Log.Debugf(MSG_FILE_DOWLOADED, fileName)
			// Juntar todas las partes del archivo
			blocks.RestoreFile(fileName)
		}
	}
	return nil
}

// Find block localmente. Retorna <endFile><nextBlockKey><error>. En caso de no poder enviar el mensaje retorna error
func (node *Node) findBlockLocally(key []byte) (bool, []byte, error) {
	fileNameFound, data, err := node.KeyValueTab.Get(key)
	// si no se encuentra retorna error
	if err != nil {
		return false, helpers.GetNullKey(), err
	}
	// si se encuentra guarda localmente y parsear data
	endFile, _ := file_manager.StoreBlockOnDownload(fileNameFound, data)
	return endFile, blocks.GetNextBlock(data), nil
}

// Se encarga de tomar contactos de un canal y hacer un ping a dichos contactos para luego
// agregarlos a la tabla
func (node *Node) PendingPingsService() {
	closed := false
	for !closed {
		contact, ok := <-node.PendingContactsToAdd
		if ok {
			if node.SndPing(node.Config, contact) == nil {
				node.AddContactPreventingLoop(contact)
			}
		}
		closed = !ok
	}
}

// Retorna una función que intenta enviar la orden de store a los vecinos más cercanos a la clave
// y en caso de no encontrar alguno almacena el bloque localmente
func (node *Node) createSndBlockNeighbors() file_manager.ProcessBlockCallBack {
	return func(key []byte, fileName string, data []byte) error {
		// buscar contactos cercanos a la clave
		contacts := node.BucketTab.GetContactsForId(key)
		// enviar mensaje de store a cada nodo encontrado
		countSended := 0
		for index := range contacts {
			if node.SndStore(node.Config, contacts[index], key, fileName, data) == nil {
				countSended++
			}
		}
		// si no se pudo agregar a algún vecino se guarda localmente
		if countSended == 0 {
			return node.doStoreBlock(key, fileName, data)
		}
		return nil
	}
}

// Agrega contactos al nodo y deja como pendiente el envío de ping's a los mismos
func (node *Node) addContactsDefferedPing(contacts []contacts_queue.Contact) {
	// Cargar contactos en el canal para hacer ping
	for _, contact := range contacts {
		node.PendingContactsToAdd <- contact
	}
	// Agrega contactos
	//node.BucketTab.AddContacts(contacts)
}

/*
	// Si no hay más contactos retornar error
	if contactStorage.IsEmpty() {
		msg := fmt.Sprintf(MSG_ERROR_FILE_NOT_FOUND, blockName)
		common.Log.Errorf(msg)
		return errors.New(msg)
	}
	// Tomar el contacto más cercano y solicitar bloque/contactos cercanos
	contact, _ := contactStorage.Pop()
	fileNameFound, nextBlockKeyFound, data, neighborContacts, err := node.SndFindBlock(node.Config, *contact, key)
	if err != nil {
		continue
	}
	// Agrego contacto a lista local
	node.AddContactPreventingLoop(*contact)
	if len(fileNameFound) > 0 {
		endFile, _ = file_manager.StoreBlockOnDownload(fileNameFound, data)
		endBlock = true
		key = nextBlockKeyFound
		common.Log.Debugf(MSG_FILE_FOUND, fileNameFound)
		continue
	}
	// Agregar contactos encontrados a la búsqueda
	if neighborContacts != nil {
		common.Log.Debugf(MSG_CONTACTS_ADDED_FOR_SEARCH, len(neighborContacts))
		contactStorage.PushContacts(neighborContacts)
	}*/
