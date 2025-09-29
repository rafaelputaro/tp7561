package protoUtils

import (
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

// Retorna el operando para la operación find block
func CreateFindBlockOperands(sourceId []byte, sourceUrl string, key []byte) *protopb.FindBlockOpers {
	return &protopb.FindBlockOpers{
		SourceId:  sourceId,
		SourceUrl: proto.String(sourceUrl),
		BlockKey:  key,
	}
}

// Retorna el resultado de la operación find block
func CreateFindBlockResults(blockName string, data []byte, contacts []contacts_queue.Contact) *protopb.FindBlockRes {
	contacstIds, contactsUrls := contactsToArrays(contacts)
	return &protopb.FindBlockRes{
		BlockName:    &blockName,
		BlockData:    data,
		ContactsIds:  contacstIds,
		ContactsUrls: contactsUrls,
	}
}

// Realiza el parseo de los operando recibidos en una operación de encontrar bloque
func ParseFindBlockOperands(operands *protopb.FindBlockOpers) (contacts_queue.Contact, []byte) {
	contactSource := contacts_queue.NewContact(operands.GetSourceId(), operands.GetSourceUrl())
	return *contactSource, operands.GetBlockKey()
}

// Pasea los resultados de una operación de encontrar bloque <fileName>,<data>,<contacts>
func ParseFindBlockResults(result *protopb.FindBlockRes) (string, []byte, []contacts_queue.Contact) {
	return result.GetBlockName(), result.GetBlockData(), contactsFromArrays(result.ContactsIds, result.ContactsUrls)
}
