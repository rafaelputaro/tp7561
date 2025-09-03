package bucket_table

// Es un contacto el cual esta dada por un id y una url
type Contact struct {
	ID  []byte
	Url string
}

// Representa una cola de contactos tipo FIFO
type ContactQueue struct {
	entries []Contact
}

// Retorna una instancia de entrada de cola lista para ser utilizada
func NewQueueContacts(id []byte, url string) *Contact {
	entry := Contact{
		ID:  id,
		Url: url,
	}
	return &entry
}

// Retorna una instancia de cola de entradas lista para ser utilizada
func NewQueue() *ContactQueue {
	queue := ContactQueue{
		entries: []Contact{},
	}
	return &queue
}

// Encola una entrada
func (queue *ContactQueue) Enqueue(entry Contact) {
	queue.entries = append(queue.entries, entry)
}

// Obtiene el top de la cola
func (queue *ContactQueue) Top() Contact {
	return queue.entries[0]
}

// Desencola el elemento top de la cola
func (queue *ContactQueue) Dequeue() Contact {
	temp := queue.entries[0]
	queue.entries = queue.entries[1:]
	return temp
}

// Retorna verdadero si la cola esta vacía
func (queue *ContactQueue) Empty() bool {
	return len(queue.entries) == 0
}

// Retorna todas las entradas de la cola
func (queue *ContactQueue) GetContacs() []Contact {
	return queue.entries
}
