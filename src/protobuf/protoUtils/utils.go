package protoUtils

import (
	"tp/peer/dht/bucket_table/contacts_queue"
)

// Convierte una lista de contactos en sendas listas de ids y de url
func contactsToArrays(contacts []contacts_queue.Contact) ([][]byte, []string) {
	contacstIds := [][]byte{}
	contactsUrls := []string{}
	for i := range contacts {
		contacstIds = append(contacstIds, contacts[i].ID)
		contactsUrls = append(contactsUrls, contacts[i].Url)
	}
	return contacstIds, contactsUrls
}

// Crea una lista de contactos en base a sendas listas de ids y urls
func contactsFromArrays(ids [][]byte, urls []string) []contacts_queue.Contact {
	contacts := []contacts_queue.Contact{}
	for i := range ids {
		contacts = append(contacts, *contacts_queue.NewContact(ids[i], urls[i]))
	}
	return contacts
}
