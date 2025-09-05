package dht

import (
	"errors"
	"testing"

	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
)

func TestNode(t *testing.T) {
	name := "peer-1"
	url := helpers.GenerateURL(name, "5001")
	id := []byte{128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	entriesPerBucket := 20
	config := helpers.PeerConfig{
		Name:              name,
		Url:               url,
		EntriesPerKBucket: entriesPerBucket,
		Id:                id,
	}
	config.LogConfig()
	node := NewNode(config, PingOpWithoutError, StoreOpWithoutError)
	// Agregar contacto durante ping
	idContact := []byte{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	node.Ping(*contacts_queue.NewContact(idContact, "contact:5001"))
	// Rechazar un ping a si mismo
	node.Ping(*contacts_queue.NewContact(config.Id, config.Url))
	//key := helpers.GetKey("")
	/*	key := []byte{}
		key = append(key, 4)
		prefixes := helpers.GeneratePrefixesOtherTrees(key)
		fmt.Println(prefixes)*/
	//arrayBool := helpers.ConvertToBoolArray(key)
	//print("%v", fmt.Sprintf("%v", arrayBool))
	//print("%v", len(arrayBool))
	/*	callback := func(url string) bool {
			return false
		}

		table := NewBucketTable(key, 10, callback)
		contact := []byte{}
		contact = append(contact, 5)
		table.AddContact(contact, "contact5::5051")
		contacts := table.GetContactsForId(contact)
		println("contact: %v", contacts[0].Url)
	*/

	/*
		contactF := table.DequeueContact("00000101")
		if contactF == nil {
			println("No hay contactos")
		} else {
			println("%v", contactF.ID)
		}*/

}

func PingOpWithError(config helpers.PeerConfig, contact contacts_queue.Contact) error {
	return errors.New("Error")
}

func PingOpWithoutError(config helpers.PeerConfig, contact contacts_queue.Contact) error {
	return nil
}

func StoreOpWithError(config helpers.PeerConfig, contact contacts_queue.Contact, key []byte, value string) error {
	return errors.New("Error")
}

func StoreOpWithoutError(config helpers.PeerConfig, contact contacts_queue.Contact, key []byte, value string) error {
	return errors.New("Error")
}
