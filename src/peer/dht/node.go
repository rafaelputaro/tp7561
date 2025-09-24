package dht

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"tp/common"
	"tp/peer/dht/bucket_table"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/dht/key_value_table"

	"tp/peer/helpers"
	"tp/peer/helpers/communication/rpc_ops"
	"tp/peer/helpers/file_manager"
	"tp/peer/helpers/file_manager/blocks"
)

const MSG_ERROR_OWN_REQUEST = "it is my own request"
const MSG_MUST_DISCARD_CONTACT = "Contact request should be discarded: %v"

// Representa un nodo de una Distributed Hash Table
type Node struct {
	Config                helpers.PeerConfig
	BucketTab             bucket_table.BucketTable
	KeyValueTab           key_value_table.KeyValueTable
	SndStore              rpc_ops.StoreOp
	SndShareContactsRecip rpc_ops.SndShareContactsRecipOp
	SndPing               rpc_ops.PingOp
	SndFindBlock          rpc_ops.FindBlockOp
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
	}
	return node
}

// Retorna verdadero si la instancia es el bootstrap node
func (node *Node) IsBootstrapNode() bool {
	return bytes.Equal(node.Config.Id, helpers.BootstrapNodeID)
}

// Chequea si es un contacto a si mismo e intenta agregarlo
func (node *Node) AddContactPreventingLoop(sourceContact contacts_queue.Contact) bool {
	// Prevenir bucle
	if node.DiscardContact(sourceContact) {
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
func (node *Node) RcvShareContactsReciprocally(sourceContact contacts_queue.Contact, sourceContactList []contacts_queue.Contact) []contacts_queue.Contact {
	// obtener contactos recomendados
	selfContacts := node.BucketTab.GetRecommendedContactsForId(sourceContact.ID)
	// agregar contactos que compartió la fuente
	node.BucketTab.AddContacts(selfContacts)
	return selfContacts
}

// Envía los contactos propios al bootstrap node esperando que el mismo retorne los contactos recomendados
// para la clave del presente nodo
func (node *Node) SndShareContactsToBootstrap() {
	if !node.IsBootstrapNode() {
		contactBoostrapNode := contacts_queue.NewContact(helpers.BootstrapNodeID, helpers.BootstrapNodeUrl)
		node.SndShareContacts(*contactBoostrapNode)
	}
}

// Envía los contactos propios al contacto node esperando que el mismo retorne los contactos recomendados
// para la clave del presente nodo
func (node *Node) SndShareContacts(destContact contacts_queue.Contact) error {
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
	node.BucketTab.AddContacts(destRcvContacts)
	return nil
}

// Retorna los contactos de los nodos más cercanos a un targetId. Además hace el intento de
// agregar el contacto solicitante a la bucket_table
func (node *Node) RcvFindNode(sourceContact contacts_queue.Contact, targetId []byte) []contacts_queue.Contact {
	// Prevenir bucle
	if node.DiscardContact(sourceContact) {
		common.Log.Debugf(fmt.Sprintf(MSG_MUST_DISCARD_CONTACT, sourceContact.ToString()))
		return []contacts_queue.Contact{}
	}
	// Agregar contacto a la bucket_table
	node.AddContactPreventingLoop(sourceContact)
	//node.BucketTab.AddContact(sourceContact)
	// Buscar los contactos
	return node.BucketTab.GetContactsForId(targetId)
}

// Si la target key se encuentra en el nodo retorna el archivo asociado a la misma, caso contrario retorna
// un error y la lista de los contactos más cercanos a la misma. Además hace el intento de
// agregar el contacto solicitante a la bucket_table
func (node *Node) RcvFindBlock(sourceContact contacts_queue.Contact, targetKey []byte) (string, []byte, []contacts_queue.Contact, error) {
	// Prevenir bucle
	if node.DiscardContact(sourceContact) {
		common.Log.Debugf(fmt.Sprintf(MSG_MUST_DISCARD_CONTACT, sourceContact.ToString()))
		return key_value_table.EMPTY_VALUE, []byte{}, []contacts_queue.Contact{}, errors.New(MSG_ERROR_OWN_REQUEST)
	}
	// Agregar contacto a la bucket_table
	node.AddContactPreventingLoop(sourceContact)
	//node.BucketTab.AddContact(sourceContact)
	// Búsqueda de archivo
	fileName, data, err := node.KeyValueTab.Get(targetKey)
	if err == nil {
		return fileName, data, []contacts_queue.Contact{}, nil
	}
	contactsToReturn := node.BucketTab.GetContactsForId(targetKey)
	return key_value_table.EMPTY_VALUE, []byte{}, contactsToReturn, err
}

// Almacena la clave valor localmente y envía el menseja de store a los contactos más cercanos a la tabla.
// En caso de que la clave ya existía localmente retorna error. Por otro lado intenta agregar el contacto
// fuente en la tabla de contactos
func (node *Node) RcvStore(sourceContact contacts_queue.Contact, key []byte, fileName string, data []byte) error {
	// Prevenir bucle
	if node.DiscardContact(sourceContact) {
		common.Log.Debugf(fmt.Sprintf(MSG_MUST_DISCARD_CONTACT, sourceContact.ToString()))
		return errors.New(MSG_ERROR_OWN_REQUEST)
	}
	// Agregar contacto a la bucket_table
	node.AddContactPreventingLoop(sourceContact)
	//node.BucketTab.AddContact(sourceContact)
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
func (node *Node) DiscardContact(contact contacts_queue.Contact) bool {
	return node.Config.Url == contact.Url
}

// Retorna los contactos para un id dado
func (node *Node) GetContactsForId(id []byte) []contacts_queue.Contact {
	return node.BucketTab.GetContactsForId(id)
}

// Obtiene el valor para una clave. En caso de no disponer la clave retorna error
func (node *Node) GetValue(key []byte) (string, []byte, error) {
	return node.KeyValueTab.Get(key)
}

// Agrega un archivo del espacio local al ipfs dado por los nodos de la red de contactos
func (node *Node) AddFile(fileName string) error {
	return file_manager.AddFile(fileName, node.createSndBlockNeighbors())
}

func (node *Node) GetFile(fileName string) error {
	name := blocks.GenerateBlockName(fileName, 1)
	key := helpers.GetKey(name)
	contacts := node.GetContactsForId(key)
	if contacts != nil {
		fileName, data, _, _ := node.SndFindBlock(node.Config, contacts[0], key)
		file_manager.StoreBlockOnDownload(fileName, data)
		common.Log.Debugf(fileName)
	} else {
		common.Log.Debugf("No hay contactos")
	}
	// obtener primer
	return nil
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
