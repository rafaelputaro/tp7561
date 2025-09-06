package dht

import (
	"bytes"
	"errors"
	"fmt"
	"tp/peer/dht/bucket_table"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
)

const BOOTSTRAP_NODE_NAME = "peer-1"
const BOOTSTRAP_NODE_PORT = "50051"
const BOOTSTRAP_NODE_HOST = "peer-1"
const MSG_ERROR_OWN_REQUEST = "it is my own request"
const MSG_MUST_DISCARD_CONTACT = "Contact request should be discarded: %v"

var BootstrapNodeUrl = helpers.GenerateURL(BOOTSTRAP_NODE_HOST, BOOTSTRAP_NODE_PORT)
var BootstrapNodeID = helpers.GetKey(BOOTSTRAP_NODE_NAME)

type StoreOp func(config helpers.PeerConfig, contact contacts_queue.Contact, key []byte, value string) error

type Node struct {
	Config      helpers.PeerConfig
	BucketTab   bucket_table.BucketTable
	KeyValueTab KeyValueTable
	SndStore    StoreOp
	// cache
}

// Retorna una nueva instancia de nodo lista para ser utilizada
func NewNode(config helpers.PeerConfig, sndPing bucket_table.PingOp, sndStore StoreOp) *Node {
	node := &Node{
		Config:      config,
		BucketTab:   *bucket_table.NewBucketTable(config, sndPing),
		KeyValueTab: *NewKeyValueTable(),
		SndStore:    sndStore,
	}
	return node
}

// Retorna verdadero si la instancia es el bootstrap node
func (node *Node) IsBootstrapNode() bool {
	return bytes.Equal(node.Config.Id, BootstrapNodeID) //node.Config.Url == BootstrapNodeUrl
}

// Simula un ping el cuál consiste en intentar agregar el contacto a la tabla de contactos
func (node *Node) Ping(contactSource contacts_queue.Contact) {
	// Prevenir bucle
	if node.DiscardContact(contactSource) {
		helpers.Log.Debugf(fmt.Sprintf(MSG_MUST_DISCARD_CONTACT, contactSource))
		return
	}
	// Trata de agregar el contacto
	node.BucketTab.AddContact(contactSource)
}

// Retorna los contactos de los nodos más cercanos a un targetId. Además hace el intento de
// agregar el contacto solicitante a la bucket_table
func (node *Node) FindNode(contactSource contacts_queue.Contact, targetId []byte) []contacts_queue.Contact {
	// Prevenir bucle
	if node.DiscardContact(contactSource) {
		helpers.Log.Debugf(fmt.Sprintf(MSG_MUST_DISCARD_CONTACT, contactSource))
		return []contacts_queue.Contact{}
	}
	// Agregar contacto a la bucket_table
	node.BucketTab.AddContact(contactSource)
	// Buscar los contactos
	return node.BucketTab.GetContactsForId(targetId)
}

// Si la target key se encuentra en el nodo retorna el valor de la misma, caso contrario retorna
// un error y la lista de los contactos más cercanos a la misma. Además hace el intento de
// agregar el contacto solicitante a la bucket_table
func (node *Node) FindValue(contactSource contacts_queue.Contact, targetKey []byte) (string, []contacts_queue.Contact, error) {
	// Prevenir bucle
	if node.DiscardContact(contactSource) {
		helpers.Log.Debugf(fmt.Sprintf(MSG_MUST_DISCARD_CONTACT, contactSource))
		return EMPTY_VALUE, []contacts_queue.Contact{}, errors.New(MSG_ERROR_OWN_REQUEST)
	}
	// Agregar contacto a la bucket_table
	node.BucketTab.AddContact(contactSource)
	// Búsqueda de valor
	valueToReturn, err := node.KeyValueTab.GetValue(targetKey)
	if err == nil {
		return valueToReturn, nil, nil
	}
	contactsToReturn := node.BucketTab.GetContactsForId(targetKey)
	return EMPTY_VALUE, contactsToReturn, err
}

// Almacena la clave valor localmente y envía el menseja de store a los contactos más cercanos a la tabla.
// En caso de que la clave ya existía localmente retorna error. Por otro lado intenta agregar el contacto
// fuente en la tabla de contactos
func (node *Node) Store(contactSource contacts_queue.Contact, key []byte, value string) error {
	// Prevenir bucle
	if node.DiscardContact(contactSource) {
		helpers.Log.Debugf(fmt.Sprintf(MSG_MUST_DISCARD_CONTACT, contactSource))
		return errors.New(MSG_ERROR_OWN_REQUEST)
	}
	// Agregar contacto a la bucket_table
	node.BucketTab.AddContact(contactSource)
	// Almacenar localmente
	err := node.KeyValueTab.Add(key, value)
	if err != nil {
		return err
	}
	// Buscar contactos cercanos
	contacts := node.BucketTab.GetContactsForId(key)
	for index := range contacts {
		node.SndStore(node.Config, contacts[index], key, value)
	}
	return nil
}

// Retorna verdadero si la url propia y la del contacto coinciden
func (node *Node) DiscardContact(contact contacts_queue.Contact) bool {
	return node.Config.Url == contact.Url
}
