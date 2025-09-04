package bucket_table

import (
	"errors"
	"slices"
	"tp/peer/helpers"
)

const MSG_ERROR_PREFIX_NOT_FOUND = "error prefix not found"
const MSG_ERROR_ON_ENQUEUE_CONTACT = "error on enqueue contact"

// Retorna verdadero si la url no se encuentra respondiendo a request's.
type IsUnresponsiveUrl func(url string) bool

// Es una tabla que contiene los contactos por prefijo
type BucketTable struct {
	Entries     map[string]ContactQueue
	Prefixes    []string
	IsUnrespUrl IsUnresponsiveUrl
}

// Retorna una tabla de contactos lista para ser utilizada
func NewBucketTable(id []byte, maxContactsPrefix int, isUnrespUrl IsUnresponsiveUrl) *BucketTable {
	table := BucketTable{
		Entries:     map[string]ContactQueue{},
		Prefixes:    []string{},
		IsUnrespUrl: isUnrespUrl,
	}
	table.initPrefixes(id)
	table.initEntries(maxContactsPrefix)
	return &table
}

// Inicializa la lista de prefijos para un id dado
func (table *BucketTable) initPrefixes(id []byte) {
	table.Prefixes = helpers.GeneratePrefixesOtherTrees(id)
}

// Inicializa las colas correspondientes a cada uno de los prefijos
func (table *BucketTable) initEntries(capacity int) {
	for prefix := range table.Prefixes {
		table.Entries[table.Prefixes[prefix]] = *NewQueue(capacity)
	}
}

// Si la tabla no se encuentra llena
func (table *BucketTable) AddContact(id []byte, url string) error {
	prefix, err := table.getPrefix(id)
	if err == nil {
		newContact := NewContact(id, url)
		queue := table.Entries[prefix]
		err := queue.Enqueue(*newContact)
		if err != nil {
			headContact, _ := queue.TakeHead()
			if table.isUnresponsiveContact(headContact) {
				queue.Enqueue(*newContact)
			} else {
				queue.Enqueue(headContact)
			}
		}
		table.Entries[prefix] = queue
		return nil
	}
	return errors.New(MSG_ERROR_ON_ENQUEUE_CONTACT)
}

// Retorna verdadero si el contacto no se encuentra resposivo
func (table *BucketTable) isUnresponsiveContact(contact Contact) bool {
	return table.IsUnrespUrl(contact.Url)
}

// Obtiene todos los contactos cercanos a un id dado
func (table *BucketTable) GetContactsForId(id []byte) []Contact {
	prefix, error := table.getPrefix(id)
	if error != nil {
		return []Contact{}
	}
	toReturn := table.GetContactsForPrefix(prefix)
	return toReturn
}

// Obtiene los contactos para un prefijo dado
func (table *BucketTable) GetContactsForPrefix(prefix string) []Contact {
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
