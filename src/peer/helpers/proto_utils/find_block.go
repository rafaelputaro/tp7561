package proto_utils_peer

import (
	"tp/common/contact"
	"tp/peer/helpers/file_manager/blocks"
	"tp/protobuf/protoUtils"
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
func CreateFindBlockResults(blockName string, data []byte, contacts []contact.Contact) *protopb.FindBlockRes {
	contacstIds, contactsUrls := protoUtils.ContactsToArrays(contacts)
	return &protopb.FindBlockRes{
		BlockName:    &blockName,
		BlockData:    data,
		ContactsIds:  contacstIds,
		ContactsUrls: contactsUrls,
	}
}

// Realiza el parseo de los operando recibidos en una operación de encontrar bloque
func ParseFindBlockOperands(operands *protopb.FindBlockOpers) (contact.Contact, []byte) {
	contactSource := contact.NewContact(operands.GetSourceId(), operands.GetSourceUrl())
	return *contactSource, operands.GetBlockKey()
}

// Pasea los resultados de una operación de encontrar bloque <fileName>,<next block key>,<data>,<contacts>
func ParseFindBlockResults(result *protopb.FindBlockRes) (string, []byte, []byte, []contact.Contact) {
	return result.GetBlockName(), blocks.GetNextBlock(result.BlockData), result.GetBlockData(), protoUtils.ContactsFromArrays(result.ContactsIds, result.ContactsUrls)
}
