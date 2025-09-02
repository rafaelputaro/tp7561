package bucket_table

type BucketTableEntry struct {
	Prefix string
	Queue  EntriesQueue
}

// Es una table que contiene pares clave valor
type TBucketTable struct {
	Entries map[string]string
}
