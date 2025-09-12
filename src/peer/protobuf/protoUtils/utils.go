package protoUtils

import (
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

// Retorna los operandos para hacer una operación de ping
func CreatePingOperands(id []byte, url string) *protopb.PingOperands {
	return &protopb.PingOperands{
		SourceId:  id,
		SourceUrl: proto.String(url),
	}
}

// Retorna los operandos para hacer compartir contactos
func CreateShareContactsReciprocallyOperands(contact contacts_queue.Contact, contacts []contacts_queue.Contact) *protopb.ShareContactsReciprocallyOperands {
	contacstIds, contactsUrls := contactsToArrays(contacts)
	return &protopb.ShareContactsReciprocallyOperands{
		SourceId:     contact.ID,
		SourceUrl:    proto.String(contact.Url),
		ContactsIds:  contacstIds,
		ContactsUrls: contactsUrls,
	}
}

// Hace el parseo de los operando recibidos en una operación de compartir contactos
func ParseShareContactsReciprocallyOperands(operands *protopb.ShareContactsReciprocallyOperands) (contacts_queue.Contact, []contacts_queue.Contact) {
	contactSource := contacts_queue.NewContact(operands.GetSourceId(), operands.GetSourceUrl())
	contacts := contactsFromArrays(operands.GetContactsIds(), operands.GetContactsUrls())
	return *contactSource, contacts
}

// Crea los resultados ha retorna en una operación de compartir contactos
func CreateShareContactsReciprocallyResults(contacts []contacts_queue.Contact) *protopb.ShareContactsReciprocallyResults {
	contacstIds, contactsUrls := contactsToArrays(contacts)
	return &protopb.ShareContactsReciprocallyResults{
		ContactsIds:  contacstIds,
		ContactsUrls: contactsUrls,
	}
}

// Pasea los resultados de una operación de compartir contactos
func ParseShareContactsReciprocallyResults(result *protopb.ShareContactsReciprocallyResults) []contacts_queue.Contact {
	toReturn := []contacts_queue.Contact{}
	for i := range result.ContactsIds {
		toReturn = append(toReturn, *contacts_queue.NewContact(result.ContactsIds[i], result.ContactsUrls[i]))
	}
	return toReturn
}

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
