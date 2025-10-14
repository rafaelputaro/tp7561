package protoUtils

import (
	"tp/common/contact"
	"tp/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

// Retorna los operandos para hacer compartir contactos
func CreateShareContactsReciprocallyOperands(contact contact.Contact, contacts []contact.Contact) *protopb.ShCtsRecipOpers {
	contacstIds, contactsUrls := ContactsToArrays(contacts)
	return &protopb.ShCtsRecipOpers{
		SourceId:     contact.ID,
		SourceUrl:    proto.String(contact.Url),
		ContactsIds:  contacstIds,
		ContactsUrls: contactsUrls,
	}
}

// Crea los resultados ha retorna en una operación de compartir contactos
func CreateShareContactsReciprocallyResults(contacts []contact.Contact) *protopb.ShCtsRecipRes {
	contacstIds, contactsUrls := ContactsToArrays(contacts)
	return &protopb.ShCtsRecipRes{
		ContactsIds:  contacstIds,
		ContactsUrls: contactsUrls,
	}
}

// Hace el parseo de los operando recibidos en una operación de compartir contactos
func ParseShareContactsReciprocallyOperands(operands *protopb.ShCtsRecipOpers) (contact.Contact, []contact.Contact) {
	contactSource := contact.NewContact(operands.GetSourceId(), operands.GetSourceUrl())
	contacts := ContactsFromArrays(operands.GetContactsIds(), operands.GetContactsUrls())
	return *contactSource, contacts
}

// Pasea los resultados de una operación de compartir contactos
func ParseShareContactsReciprocallyResults(result *protopb.ShCtsRecipRes) []contact.Contact {
	toReturn := []contact.Contact{}
	for i := range result.ContactsIds {
		toReturn = append(toReturn, *contact.NewContact(result.ContactsIds[i], result.ContactsUrls[i]))
	}
	return toReturn
}
