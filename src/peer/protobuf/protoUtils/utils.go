package protoUtils

import (
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

func CreatePingOperands(id []byte, url string) *protopb.PingOperands {
	return &protopb.PingOperands{
		SourceId:  id,
		SourceUrl: proto.String(url),
	}
}

func CreateShareContactsReciprocallyOperands(contact contacts_queue.Contact, contacts []contacts_queue.Contact) *protopb.ShareContactsReciprocallyOperands {
	contacstIds, contactsUrls := contactsToArrays(contacts)
	return &protopb.ShareContactsReciprocallyOperands{
		SourceId:     contact.ID,
		SourceUrl:    proto.String(contact.Url),
		ContactsIds:  contacstIds,
		ContactsUrls: contactsUrls,
	}
}

func ParseShareContactsReciprocallyOperands(operands *protopb.ShareContactsReciprocallyOperands) (contacts_queue.Contact, []contacts_queue.Contact) {
	contactSource := contacts_queue.NewContact(operands.GetSourceId(), operands.GetSourceUrl())
	contacts := contactsFromArrays(operands.GetContactsIds(), operands.GetContactsUrls())
	return *contactSource, contacts
}

func CreateShareContactsReciprocallyResults(contacts []contacts_queue.Contact) *protopb.ShareContactsReciprocallyResults {
	contacstIds, contactsUrls := contactsToArrays(contacts)
	return &protopb.ShareContactsReciprocallyResults{
		ContactsIds:  contacstIds,
		ContactsUrls: contactsUrls,
	}
}

func ParseShareContactsReciprocallyResults(result *protopb.ShareContactsReciprocallyResults) []contacts_queue.Contact {
	toReturn := []contacts_queue.Contact{}
	for i := range result.ContactsIds {
		toReturn = append(toReturn, *contacts_queue.NewContact(result.ContactsIds[i], result.ContactsUrls[i]))
	}
	return toReturn
}

func contactsToArrays(contacts []contacts_queue.Contact) ([][]byte, []string) {
	contacstIds := [][]byte{}
	contactsUrls := []string{}
	for i := range contacts {
		contacstIds = append(contacstIds, contacts[i].ID)
		contactsUrls = append(contactsUrls, contacts[i].Url)
	}
	return contacstIds, contactsUrls
}

func contactsFromArrays(ids [][]byte, urls []string) []contacts_queue.Contact {
	contacts := []contacts_queue.Contact{}
	for i := range ids {
		contacts = append(contacts, *contacts_queue.NewContact(ids[i], urls[i]))
	}
	return contacts
}
