package protoUtils

import (
	"tp/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

// Retorna el operando para la operación add file
func CreateAddFileOperands(fileName string) *protopb.AddFileOpers {
	return &protopb.AddFileOpers{
		FileName: proto.String(fileName),
	}
}

// Retorna el resultado de la operación add file
func CreateAddFileResults(key []byte, url string) *protopb.AddFileRes {
	return &protopb.AddFileRes{
		Key: key,
		Url: proto.String(url),
	}
}

// Realiza el parseo de los operando recibidos en una operación de agregar archivo
// <fileName><part><data><endFile>
func ParseAddFileOperands(operands *protopb.AddFileOpers) string {
	return operands.GetFileName()
}

// Pasea los resultados de una operación de agregar archivo <key>
func ParseAddFileResults(result *protopb.AddFileRes) ([]byte, string) {
	return result.GetKey(), result.GetUrl()
}
