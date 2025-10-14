package tiered_contact_storage

import "tp/common/contact"

// Es una pila de contactos con control de repetidos
type ContactStack struct {
	contactList []contact.Contact
	presentList map[string]bool
}

// Retorna una instancia lista para ser utilizada
func NewContactStack() *ContactStack {
	return &ContactStack{
		contactList: []contact.Contact{},
		presentList: map[string]bool{},
	}
}

// Retorna verdadero si la pila se encuentra vacía
func (stack *ContactStack) IsEmpty() bool {
	return len(stack.contactList) == 0
}

// Apila un elemento en la pila. Retorna verdadero si pudo agregarse
func (stack *ContactStack) Push(contact contact.Contact) bool {
	if _, ok := stack.presentList[contact.Url]; !ok {
		stack.presentList[contact.Url] = true
		stack.contactList = append(stack.contactList, contact)
		return true
	}
	return false
}

// Desapila el último elemento de la pila. Retorna nulo si no hay elementos en la pila
func (stack *ContactStack) Pop() *contact.Contact {
	if stack.IsEmpty() {
		return nil
	}
	toReturn := stack.contactList[len(stack.contactList)-1]
	stack.contactList = stack.contactList[:len(stack.contactList)-1]
	delete(stack.presentList, toReturn.Url)
	return &toReturn
}
