package protoUtils

import (
	"tp/peer/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

func CreatePingOperands(id []byte, url string) *protopb.PingOperands {
	return &protopb.PingOperands{
		SourceId:  id,
		SourceUrl: proto.String(url),
	}
}
