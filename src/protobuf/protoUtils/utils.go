package protoUtils

import "tp/common/contact"

// Convierte una lista de contactos en sendas listas de ids y de url
func ContactsToArrays(contacts []contact.Contact) ([][]byte, []string) {
	contacstIds := [][]byte{}
	contactsUrls := []string{}
	for i := range contacts {
		contacstIds = append(contacstIds, contacts[i].ID)
		contactsUrls = append(contactsUrls, contacts[i].Url)
	}
	return contacstIds, contactsUrls
}

// Crea una lista de contactos en base a sendas listas de ids y urls
func ContactsFromArrays(ids [][]byte, urls []string) []contact.Contact {
	contacts := []contact.Contact{}
	for i := range ids {
		contacts = append(contacts, *contact.NewContact(ids[i], urls[i]))
	}
	return contacts
}
