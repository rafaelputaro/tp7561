package bucket_table

// Es una entrada de la cola la cual esta dada por un id y una url
type QueueEntry struct {
	ID  []byte
	Url string
}

// Representa una cola de entrada tipo FIFO
type EntriesQueue struct {
	entries []QueueEntry
}

// Retorna una instancia de cola de entradas lista para ser utilizada
func NewQueue() *EntriesQueue {
	queue := EntriesQueue{
		entries: []QueueEntry{},
	}
	return &queue
}

// Encola una entrada
func (queue *EntriesQueue) Enqueue(entry QueueEntry) {
	queue.entries = append(queue.entries, entry)
}

// Obtiene el top de la cola
func (queue *EntriesQueue) Top() QueueEntry {
	return queue.entries[0]
}

// Desencola el elemento top de la cola
func (queue *EntriesQueue) Dequeue() QueueEntry {
	temp := queue.entries[0]
	queue.entries = queue.entries[1:]
	return temp
}

// Retorna verdadero si la cola esta vacía
func (queue *EntriesQueue) Empty() bool {
	return len(queue.entries) == 0
}

// Retorna todas las entradas de la cola
func (queue *EntriesQueue) GetEntries() []QueueEntry {
	return queue.entries
}
