package bucket_table

import (
	"errors"
	"fmt"
	"slices"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
)

const MSG_ERROR_PREFIX_NOT_FOUND = "error prefix not found"
const MSG_ERROR_ON_ENQUEUE_CONTACT = "error on enqueue contact"
const MSG_CONTACT_ADDED = "The contact has been added | url: %v"
const MSG_TRY_TO_ADD_CONTACTS = "Attempt to add %v contacts"
const MSG_CONTACT_REPLACE_HEAD = "Contact (url: %v) has been added to replace tailhead (url: %v)"
const MSG_CONTACT_DISCARD = "Contact has been ruled out | url: %v"

type PingOp func(config helpers.PeerConfig, contact contacts_queue.Contact) error

// Es una tabla que contiene los contactos por prefijo
type BucketTable struct {
	Config   helpers.PeerConfig
	Entries  map[string]contacts_queue.ContactQueue
	Prefixes []string
	Ping     PingOp
}

// Retorna una tabla de contactos lista para ser utilizada
func NewBucketTable(config helpers.PeerConfig, ping PingOp) *BucketTable {
	table := BucketTable{
		Config:   config,
		Entries:  map[string]contacts_queue.ContactQueue{},
		Prefixes: []string{},
		Ping:     ping,
	}
	table.initPrefixes(table.Config.Id)
	table.initEntries(table.Config.EntriesPerKBucket)
	return &table
}

// Inicializa la lista de prefijos para un id dado
func (table *BucketTable) initPrefixes(id []byte) {
	table.Prefixes = helpers.GeneratePrefixesOtherTreesAsStrings(id)
}

// Inicializa las colas correspondientes a cada uno de los prefijos
func (table *BucketTable) initEntries(capacity int) {
	for prefix := range table.Prefixes {
		table.Entries[table.Prefixes[prefix]] = *contacts_queue.NewQueue(capacity)
	}
}

// Si la tabla no se encuentra llena agrega el contacto
func (table *BucketTable) AddContact(newContact contacts_queue.Contact) error {
	prefix, err := table.getPrefix(newContact.ID)
	if err == nil {
		queue := table.Entries[prefix]
		okEnqueue, err := queue.Enqueue(newContact)
		if err != nil {
			headContact, _ := queue.TakeHead()
			if table.isUnresponsiveContact(headContact) {
				helpers.Log.Debugf(fmt.Sprintf(MSG_CONTACT_REPLACE_HEAD, newContact.ToString(), headContact.ToString()))
				queue.Enqueue(newContact)
			} else {
				helpers.Log.Debugf(fmt.Sprintf(MSG_CONTACT_DISCARD, newContact.ToString()))
				queue.Enqueue(headContact)
			}
		} else {
			if okEnqueue {
				helpers.Log.Debugf(fmt.Sprintf(MSG_CONTACT_ADDED, newContact.ToString()))
			}
		}
		table.Entries[prefix] = queue
		return nil
	}
	helpers.Log.Errorf(MSG_ERROR_ON_ENQUEUE_CONTACT)
	return errors.New(MSG_ERROR_ON_ENQUEUE_CONTACT)
}

// Intenta agregar los contactos según la capacidad actual de la tabla
func (table *BucketTable) AddContacts(newContacts []contacts_queue.Contact) error {
	helpers.Log.Debugf(MSG_TRY_TO_ADD_CONTACTS, len(newContacts))
	for _, contact := range newContacts {
		err := table.AddContact(contact)
		if err != nil {
			return err
		}
	}
	return nil
}

// Retorna verdadero si el contacto no se encuentra resposivo
func (table *BucketTable) isUnresponsiveContact(contact contacts_queue.Contact) bool {
	err := table.Ping(table.Config, contact)
	return err != nil
}

// Selecciona de la tabla de contactos propias una serie de contactos recomendados para que
// el nodo con el id parámetro pueda armar su tabla de contactos
func (table *BucketTable) GetRecommendedContactsForId(id []byte) []contacts_queue.Contact {
	prefixes := helpers.GenerateKeysFromOtherTrees(id)
	toReturn := []contacts_queue.Contact{}
	idsMap := map[string]bool{}
	for i := range prefixes {
		contactsPref := table.GetContactsForId(prefixes[i])
		for _, contact := range contactsPref {
			idStr := helpers.KeyToString(contact.ID)
			if idsMap[idStr] {
				continue
			}
			idsMap[idStr] = true
			toReturn = append(toReturn, contact)
		}
	}
	// agregar el nodo actual a la lista
	toReturn = append(toReturn, *contacts_queue.NewContact(table.Config.Id, table.Config.Url))
	return toReturn
}

// Obtiene todos los contactos cercanos a un id dado
func (table *BucketTable) GetContactsForId(id []byte) []contacts_queue.Contact {
	prefix, error := table.getPrefix(id)
	if error != nil {
		return []contacts_queue.Contact{}
	}
	toReturn := table.GetContactsForPrefix(prefix)
	return toReturn
}

// Obtiene los contactos para un prefijo dado
func (table *BucketTable) GetContactsForPrefix(prefix string) []contacts_queue.Contact {
	entries := table.Entries[prefix]
	return entries.GetContacs()
}

// Obtiene el prefijo más cercano a una clave dada
func (table *BucketTable) getPrefix(key []byte) (string, error) {
	// generar prefijos para la clave
	prefixes := helpers.GeneratePrefixes(key)
	// buscar en mi lista de prefijos
	for iPref := range prefixes {
		prefix := prefixes[iPref]
		if slices.Contains(table.Prefixes, prefix) {
			return prefix, nil
		}
	}
	return helpers.EMPTY_KEY, errors.New(MSG_ERROR_PREFIX_NOT_FOUND)
}
