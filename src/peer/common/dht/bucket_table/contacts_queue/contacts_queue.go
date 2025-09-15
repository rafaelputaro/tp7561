package contacts_queue

import (
	"errors"
	"tp/peer/common/helpers"
)

const MSG_ERROR_EMPTY_QUEUE = "error queue empty"
const MSG_ERROR_FULL_QUEUE = "error queue is full"
const MSG_CONTACT_ALREADY_ADDED = "The contact has already been added previously: %v"

// Representa una cola de contactos tipo FIFO que no permite id's repetidos
type ContactQueue struct {
	Entries       []Contact
	IdsInTheQueue map[string]bool
	Capacity      int
}

// Retorna una instancia de cola de entradas lista para ser utilizada
func NewQueue(capacity int) *ContactQueue {
	queue := ContactQueue{
		Entries:       []Contact{},
		Capacity:      capacity,
		IdsInTheQueue: map[string]bool{},
	}
	return &queue
}

// Encola una entrada si la misma no se encuentra en la cola retornando error nulo.
// En caso de que la entrada se encuentra en la cola retorna falso y error nulo.
// Si la cosa esta llena retorna falso y un error
func (queue *ContactQueue) Enqueue(entry Contact) (bool, error) {
	if !queue.hasId(entry.ID) {
		if queue.Full() {
			return false, errors.New(MSG_ERROR_FULL_QUEUE)
		}
		queue.Entries = append(queue.Entries, entry)
		queue.IdsInTheQueue[helpers.KeyToString(entry.ID)] = true
		return true, nil
	}
	helpers.Log.Debugf(MSG_CONTACT_ALREADY_ADDED, entry.ToString())
	return false, nil
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
	delete(queue.IdsInTheQueue, helpers.KeyToString(temp.ID))
	return temp, nil
}

// Quita de la cola el último elemento encolado y lo retorna
func (queue *ContactQueue) TakeHead() (Contact, error) {
	if queue.Empty() {
		return *CreateInvalidContact(), errors.New(MSG_ERROR_EMPTY_QUEUE)
	}
	toReturn := queue.Entries[len(queue.Entries)-1]
	queue.Entries = queue.Entries[:len(queue.Entries)-1]
	delete(queue.IdsInTheQueue, helpers.KeyToString(toReturn.ID))
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

// Retorna verdadero si encuentra el id en el mapa de id's presentes en la cola
func (queue *ContactQueue) hasId(id []byte) bool {
	return queue.IdsInTheQueue[helpers.KeyToString(id)]
}
