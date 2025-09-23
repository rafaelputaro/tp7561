package protoUtils

import (
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

// Retorna el operando para la operación find block
func CreateFindBlockOperands(sourceId []byte, sourceUrl string, key []byte) *protopb.FindBlockOperands {
	return &protopb.FindBlockOperands{
		SourceId:  sourceId,
		SourceUrl: proto.String(sourceUrl),
		BlockKey:  key,
	}
}

// Retorna el resultado de la operación find block
func CreateFindBlockResults(blockName string, data []byte, contacts []contacts_queue.Contact) *protopb.FindBlockResults {
	contacstIds, contactsUrls := contactsToArrays(contacts)
	return &protopb.FindBlockResults{
		BlockName:    &blockName,
		BlockData:    data,
		ContactsIds:  contacstIds,
		ContactsUrls: contactsUrls,
	}
}

// Realiza el parseo de los operando recibidos en una operación de encontrar bloque
func ParseFindBlockOperands(operands *protopb.FindBlockOperands) (contacts_queue.Contact, []byte) {
	contactSource := contacts_queue.NewContact(operands.GetSourceId(), operands.GetSourceUrl())
	return *contactSource, operands.GetBlockKey()
}

// Pasea los resultados de una operación de encontrar bloque
func ParseFindBlockResults(result *protopb.FindBlockResults) (string, []byte, []contacts_queue.Contact) {
	return result.GetBlockName(), result.GetBlockData(), contactsFromArrays(result.ContactsIds, result.ContactsUrls)
}
