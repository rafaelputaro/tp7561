package protoUtils

import (
	"tp/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

// Retorna los operandos para hacer una operación de ping
func CreatePingOperands(id []byte, url string) *protopb.PingOperands {
	return &protopb.PingOperands{
		SourceId:  id,
		SourceUrl: proto.String(url),
	}
}
