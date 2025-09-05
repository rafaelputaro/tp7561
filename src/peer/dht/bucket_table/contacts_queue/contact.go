package contacts_queue

// Es un contacto el cual esta dada por un id y una url
type Contact struct {
	ID  []byte
	Url string
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
