package protoUtils

import (
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

// Retorna los operandos para hacer compartir contactos
func CreateShareContactsReciprocallyOperands(contact contacts_queue.Contact, contacts []contacts_queue.Contact) *protopb.ShCtsRecipOpers {
	contacstIds, contactsUrls := contactsToArrays(contacts)
	return &protopb.ShCtsRecipOpers{
		SourceId:     contact.ID,
		SourceUrl:    proto.String(contact.Url),
		ContactsIds:  contacstIds,
		ContactsUrls: contactsUrls,
	}
}

// Crea los resultados ha retorna en una operación de compartir contactos
func CreateShareContactsReciprocallyResults(contacts []contacts_queue.Contact) *protopb.ShCtsRecipRes {
	contacstIds, contactsUrls := contactsToArrays(contacts)
	return &protopb.ShCtsRecipRes{
		ContactsIds:  contacstIds,
		ContactsUrls: contactsUrls,
	}
}

// Hace el parseo de los operando recibidos en una operación de compartir contactos
func ParseShareContactsReciprocallyOperands(operands *protopb.ShCtsRecipOpers) (contacts_queue.Contact, []contacts_queue.Contact) {
	contactSource := contacts_queue.NewContact(operands.GetSourceId(), operands.GetSourceUrl())
	contacts := contactsFromArrays(operands.GetContactsIds(), operands.GetContactsUrls())
	return *contactSource, contacts
}

// Pasea los resultados de una operación de compartir contactos
func ParseShareContactsReciprocallyResults(result *protopb.ShCtsRecipRes) []contacts_queue.Contact {
	toReturn := []contacts_queue.Contact{}
	for i := range result.ContactsIds {
		toReturn = append(toReturn, *contacts_queue.NewContact(result.ContactsIds[i], result.ContactsUrls[i]))
	}
	return toReturn
}
