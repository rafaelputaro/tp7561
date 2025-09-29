package protoUtils

import (
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

// Crea los operandos para realizar una operación de rpc de guardado de bloque
func CreateStoreBlockOperands(sourceId []byte, sourceUrl string, key []byte, blockName string, data []byte) *protopb.StoreBlockOpers {
	return &protopb.StoreBlockOpers{
		SourceId:  sourceId,
		SourceUrl: proto.String(sourceUrl),
		Key:       key,
		BlockName: proto.String(blockName),
		Data:      data,
	}
}

// Parsea los operandos de la operación store block a <source contact><block key><block name><data>
func ParseStoreBlockOperands(operands *protopb.StoreBlockOpers) (*contacts_queue.Contact, []byte, string, []byte) {
	sourceContact := contacts_queue.NewContact(operands.GetSourceId(), operands.GetSourceUrl())
	blockKey := operands.GetKey()
	blockName := operands.GetBlockName()
	data := operands.GetData()
	return sourceContact, blockKey, blockName, data
}
