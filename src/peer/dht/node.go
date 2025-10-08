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
	"tp/peer/helpers/task_scheduler"
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
	TaskScheduler         task_scheduler.TaskScheduler
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
		TaskScheduler:         *task_scheduler.NewTaskScheduler(),
	}
	return node
}

// Se eliminan recursos asociados
func (node *Node) DisposeNode() {
	node.TaskScheduler.DisposeTaskScheduler()
}

// Retorna verdadero si la instancia es el bootstrap node
func (node *Node) IsBootstrapNode() bool {
	return bytes.Equal(node.Config.Id, helpers.BootstrapNodeID)
}

// Representa la recepción de un ping el cuál consiste en intentar agregar el contacto a la tabla de
// contactos
func (node *Node) RcvPing(sourceContact contacts_queue.Contact) bool {
	return node.BucketTab.AddContact(sourceContact) == nil
}

// Obtiene los contactos locales recomendados para la fuente, agrega los contactos compartidos por la fuente y
// retorna los contactos recomendados para la fuente
func (node *Node) RcvShCtsRecip(sourceContact contacts_queue.Contact, sourceContactList []contacts_queue.Contact) []contacts_queue.Contact {
	// agregar contacto origen
	node.scheduleAddContactTask(sourceContact)
	// obtener contactos recomendados
	newContacts := node.BucketTab.GetRecommendedContactsForId(sourceContact.ID)
	// agregar contactos que compartió la fuente
	node.scheduleAddContactsTask(sourceContactList)
	return newContacts
}

// Envía los contactos propios al bootstrap node esperando que el mismo retorne los contactos recomendados
// para la clave del presente nodo
func (node *Node) SndShCtsToBootstrap() {
	if !node.IsBootstrapNode() {
		contactBoostrapNode := contacts_queue.NewContact(helpers.BootstrapNodeID, helpers.BootstrapNodeUrl)
		// agregar bootstrap node a contactos
		if node.SndShCts(*contactBoostrapNode) == nil {
			node.scheduleAddContactTask(*contactBoostrapNode)
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
	node.schedulePingAndAddContactsTask(destRcvContacts)
	return nil
}

// Retorna los contactos de los nodos más cercanos a un targetId. Además hace el intento de
// agregar el contacto solicitante a la bucket_table
func (node *Node) RcvFindNode(sourceContact contacts_queue.Contact, targetId []byte) []contacts_queue.Contact {
	node.scheduleAddContactTask(sourceContact)
	// Buscar los contactos
	return node.BucketTab.GetContactsForId(targetId)
}

// Si la target key se encuentra en el nodo retorna el archivo asociado a la misma, caso contrario retorna
// un error y la lista de los contactos más cercanos a la misma. Además hace el intento de
// agregar el contacto solicitante a la bucket_table
func (node *Node) RcvFindBlock(sourceContact contacts_queue.Contact, targetKey []byte) (string, []byte, []contacts_queue.Contact, error) {
	node.scheduleAddContactTask(sourceContact)
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
	node.scheduleAddContactTask(sourceContact)
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
	node.scheduleSndStoreTask(key, fileName, data, contacts)
	return nil
}

// Agrega la tarea de envío de SndStore a un lote de contactos
func (node *Node) scheduleSndStoreTask(key []byte, fileName string, data []byte, contacts []contacts_queue.Contact) {
	for _, contact := range contacts {
		node.TaskScheduler.AddTask(func() {
			node.SndStore(node.Config, contact, key, fileName, data)
		})
	}
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
		} else {
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
				resChan := make(chan processNextContactReturn, node.Config.SearchWorkers)
				defer close(resChan)
				wg := new(sync.WaitGroup)
				// Procesar contacto
				for id := range node.Config.SearchWorkers {
					wg.Add(1)
					go func() {
						processNextContact(node, key, fileName, contactStorage, resChan, id)
						wg.Done()
					}()
				}
				wg.Wait()
				for range node.Config.SearchWorkers {
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
		}
		if endFile {
			common.Log.Debugf(MSG_FILE_DOWLOADED, fileName)
			// Juntar todas las partes del archivo
			blocks.RestoreFile(fileName)
		}
	}
	return nil
}

// Busca el bloque localmente. Retorna <endFile><nextBlockKey><error>. En caso de no poder
// enviar el mensaje retorna error
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

// Agrega la tarea de agregar un contacto a la bucket table. Se recomienda utilizarla
// para evitar posibles retrasos durante la actualización de la bucket table que impliquen
// enviar pings secundarios a otros contactos
func (node *Node) scheduleAddContactTask(contact contacts_queue.Contact) {
	node.TaskScheduler.AddTask(func() {
		node.BucketTab.AddContact(contact)
	})
}

// Agrega la tarea de agregar varios contactos a la bucket table. Se recomienda utilizarla
// para evitar posibles retrasos durante la actualización de la bucket table que impliquen
// enviar pings secundarios a otros contactos
func (node *Node) scheduleAddContactsTask(contacts []contacts_queue.Contact) {
	node.TaskScheduler.AddTask(func() {
		node.BucketTab.AddContacts(contacts)
	})
}

// Agrega la tarea de enviar ping a contactos para ser agregador a la bucket table
// en caso de encontrarse activos
func (node *Node) schedulePingAndAddContactsTask(contacts []contacts_queue.Contact) {
	for _, contact := range contacts {
		node.TaskScheduler.AddTask(func() {
			if node.SndPing(node.Config, contact) == nil {
				node.BucketTab.AddContact(contact)
			}
		})
	}
}
