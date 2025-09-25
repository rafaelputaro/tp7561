package tiered_contact_storage

import (
	"tp/peer/dht/bucket_table/contacts_queue"
)

// Es una pila de contactos
type ContactStack struct {
	Contacts []contacts_queue.Contact
}

// Retorna una instancia lista para ser utilizada
func NewContactStack() *ContactStack {
	return &ContactStack{
		Contacts: []contacts_queue.Contact{},
	}
}

// Retorna verdadero si la pila se encuentra vacía
func (stack *ContactStack) IsEmpty() bool {
	return len(stack.Contacts) == 0
}

// Apila un elemento en la pila
func (stack *ContactStack) Push(contact contacts_queue.Contact) {
	stack.Contacts = append(stack.Contacts, contact)
}

// Desapila el último elemento de la pila. Retorna nulo si no hay elementos en la pila
func (stack *ContactStack) Pop() *contacts_queue.Contact {
	if stack.IsEmpty() {
		return nil
	}
	toReturn := stack.Contacts[len(stack.Contacts)-1]
	stack.Contacts = stack.Contacts[:len(stack.Contacts)-1]
	return &toReturn
}
