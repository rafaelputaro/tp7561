package protoUtils

import (
	"tp/protobuf/protopb"

	"google.golang.org/protobuf/proto"
)

// Retorna el operando para la operación find block
func CreateAddFileOperands(fileName string, part int32, data []byte, endFile bool) *protopb.AddFileOpers {

	return &protopb.AddFileOpers{
		FileName: proto.String(fileName),
		Part:     proto.Int32(part),
		Data:     data,
		Endfile:  proto.Bool(endFile),
	}
}

// Retorna el resultado de la operación find block
func CreateAddFileResults(key []byte) *protopb.AddFileRes {
	return &protopb.AddFileRes{
		Key: key,
	}
}

// Realiza el parseo de los operando recibidos en una operación de agregar archivo
// <fileName><part><data><endFile>
func ParseAddFileOperands(operands *protopb.AddFileOpers) (string, int32, []byte, bool) {
	return operands.GetFileName(), operands.GetPart(), operands.GetData(), operands.GetEndfile()
}

// Pasea los resultados de una operación de agregar archivo <key>
func ParseAddFileResults(result *protopb.AddFileRes) []byte {
	return result.GetKey()
}
