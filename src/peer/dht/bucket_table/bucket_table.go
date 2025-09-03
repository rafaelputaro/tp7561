package bucket_table

import (
	"errors"
	"slices"
	"tp/peer/helpers"
)

const MSG_ERROR_PREFIX_NOT_FOUND = "error prefix not found"
const MSG_ERROR_ON_ENQUEUE_CONTACT = "error on enqueue contact"

// Es una table que contiene pares clave valor
type BucketTable struct {
	Entries  map[string]ContactQueue
	Prefixes []string
}

func NewBucketTable(id []byte) *BucketTable {
	table := BucketTable{
		Entries:  map[string]ContactQueue{},
		Prefixes: []string{},
	}
	table.initPrefixes(id)
	return &table
}

func (table *BucketTable) initPrefixes(id []byte) {
	table.Prefixes = helpers.GeneratePrefixesOtherTrees(id)
}

func (table *BucketTable) EnqueueContact(id []byte, url string) error {
	prefix, err := table.getPrefix(id)
	if err == nil {
		entry := NewQueueContacts(id, url)
		queue := table.Entries[prefix]
		queue.Enqueue(*entry)
		table.Entries[prefix] = queue
		return nil
	}
	return errors.New(MSG_ERROR_ON_ENQUEUE_CONTACT)
}

func (table *BucketTable) DequeueContact(prefix string) *Contact {
	queue := table.Entries[prefix]
	if queue.Empty() {
		return nil
	}
	toReturn := queue.Dequeue()
	table.Entries[prefix] = queue
	return &toReturn
}

func (table *BucketTable) GetContactsForId(id []byte) []Contact {
	prefix, error := table.getPrefix(id)
	if error != nil {
		return []Contact{}
	}
	toReturn := table.GetContactsForPrefix(prefix)
	return toReturn
}

func (table *BucketTable) GetContactsForPrefix(prefix string) []Contact {
	entries := table.Entries[prefix]
	return entries.GetContacs()
}

func (table *BucketTable) getPrefix(key []byte) (string, error) {
	prefixes := helpers.GeneratePrefixes(key)
	for iPref := range prefixes {
		prefix := prefixes[iPref]
		if slices.Contains(table.Prefixes, prefix) {
			return prefix, nil
		}
	}
	return helpers.EMPTY_KEY, errors.New(MSG_ERROR_PREFIX_NOT_FOUND)
}
