package bucket_table

import "errors"

const MSG_ERROR_EMPTY_QUEUE = "error queue empty"
const MSG_ERROR_FULL_QUEUE = "error queue is full"

// Es un contacto el cual esta dada por un id y una url
type Contact struct {
	ID  []byte
	Url string
}

// Representa una cola de contactos tipo FIFO
type ContactQueue struct {
	Entries  []Contact
	Capacity int
}

// Retorna un contacto asumido como nulo, vacío o inválido
func CreateInvalidContact() *Contact {
	return NewContact([]byte{}, "")
}

// Retorna una instancia de contacto para ser utilizada
func NewContact(id []byte, url string) *Contact {
	entry := Contact{
		ID:  id,
		Url: url,
	}
	return &entry
}

// Retorna una instancia de cola de entradas lista para ser utilizada
func NewQueue(capacity int) *ContactQueue {
	queue := ContactQueue{
		Entries:  []Contact{},
		Capacity: capacity,
	}
	return &queue
}

// Encola una entrada
func (queue *ContactQueue) Enqueue(entry Contact) error {
	if queue.Full() {
		return errors.New(MSG_ERROR_FULL_QUEUE)
	}
	queue.Entries = append(queue.Entries, entry)
	return nil
}

// Obtiene el top de la cola
func (queue *ContactQueue) Top() (Contact, error) {
	if queue.Empty() {
		return *CreateInvalidContact(), errors.New(MSG_ERROR_EMPTY_QUEUE)
	}
	return queue.Entries[0], nil
}

// Desencola el elemento top de la cola
func (queue *ContactQueue) Dequeue() (Contact, error) {
	if queue.Empty() {
		return *CreateInvalidContact(), errors.New(MSG_ERROR_EMPTY_QUEUE)
	}
	temp := queue.Entries[0]
	queue.Entries = queue.Entries[1:]
	return temp, nil
}

// Quita de la cola el último elemento encolado y lo retorna
func (queue *ContactQueue) TakeHead() (Contact, error) {
	if queue.Empty() {
		return *CreateInvalidContact(), errors.New(MSG_ERROR_EMPTY_QUEUE)
	}
	toReturn := queue.Entries[len(queue.Entries)-1]
	queue.Entries = queue.Entries[:len(queue.Entries)-1]
	return toReturn, nil
}

// Retorna verdadero si la cola esta vacía
func (queue *ContactQueue) Empty() bool {
	return len(queue.Entries) == 0
}

// Retorna verdadero si la cola esta llena
func (queue *ContactQueue) Full() bool {
	return len(queue.Entries) == queue.Capacity
}

// Retorna todas las entradas de la cola
func (queue *ContactQueue) GetContacs() []Contact {
	return queue.Entries
}
