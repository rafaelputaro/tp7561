package bucket_table

import "tp/peer/helpers"

type BucketTableEntry struct {
	Prefix string
	Queue  EntriesQueue
}

// Es una table que contiene pares clave valor
type BucketTable struct {
	Entries  map[string]BucketTableEntry
	Prefixes []string
}

func NewBucketTable(id []byte) *BucketTable {
	table := BucketTable{
		Entries:  map[string]BucketTableEntry{},
		Prefixes: []string{},
	}
	table.initPrefixes(id)
	return &table
}

func (table *BucketTable) initPrefixes(id []byte) {
	table.Prefixes = helpers.GeneratePrefixes(id)
}
