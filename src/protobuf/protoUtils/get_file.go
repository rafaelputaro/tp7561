package protoUtils

import (
	"tp/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

// Retorna el operando para la operación obtener archivo
func CreateGetFileOperands(key []byte, url string) *protopb.GetFileOpers {
	return &protopb.GetFileOpers{
		Key: key,
		Url: proto.String(url),
	}
}

// Retorna el resultado de la operación obtener archivo
func CreateGetFileResults(accepted bool, pending bool) *protopb.GetFileRes {
	return &protopb.GetFileRes{
		Accepted: proto.Bool(accepted),
		Pending:  proto.Bool(pending),
	}
}

// Realiza el parseo de los operando recibidos en una operación de obtener archivo
// <key><url>
func ParseGetFileOperands(operands *protopb.GetFileOpers) ([]byte, string) {
	return operands.GetKey(), operands.GetUrl()
}

// Pasea los resultados de una operación de obtener archivo <key>
func ParseGetFileResults(result *protopb.GetFileRes) (bool, bool) {
	return result.GetAccepted(), result.GetPending()
}
