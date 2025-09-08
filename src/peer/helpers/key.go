package helpers

import (
	"crypto/sha1"
	"encoding/hex"
)

const EMPTY_KEY = ""

// Obtiene una key SHA1 desde un string
func GetKey(data string) []byte {
	h := sha1.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}

func KeyToString(key []byte) string {
	return string(key)
}

func KeyToHexString(key []byte) string {
	hexString := hex.EncodeToString(key)
	return hexString
}

func ConvertToBinaryString(key []byte) string {
	boolArray := ConvertToBoolArray(key)
	return BoolArrayToBinaryString(boolArray)
}

func ConvertToBoolArray(data []byte) []bool {
	res := make([]bool, len(data)*8)
	for i := range res {
		res[i] = data[i/8]&(0x80>>byte(i&0x7)) != 0
	}
	return res
}

func KeyToLogFormatString(key []byte) string {
	switch LoginFormatForKeys {
	case string(HEXA):
		return KeyToHexString(key)
	default:
		return ConvertToBinaryString(key)
	}
}

// Función para convertir un array de bool en una cadena binaria
func BoolArrayToBinaryString(arr []bool) string {
	var result string
	for _, b := range arr {
		if b {
			result += "1"
		} else {
			result += "0"
		}
	}
	return result
}

// Retorna un array con los prefijos con su representación binaria como strings
func GeneratePrefixesOtherTrees(key []byte) []string {
	toReturn := []string{}
	arrayBool := ConvertToBoolArray(key)
	for k := range arrayBool {
		prefix := make([]bool, k+1)
		copy(prefix[:], arrayBool[0:k+1])
		prefix[k] = !prefix[k]
		toReturn = append(toReturn, string(BoolArrayToBinaryString(prefix)))
	}
	return toReturn
}

// Retorna un array con los prefijos con su representación binaria como strings
func GeneratePrefixes(key []byte) []string {
	toReturn := []string{}
	arrayBool := ConvertToBoolArray(key)
	for k := range arrayBool {
		prefix := make([]bool, k+1)
		copy(prefix[:], arrayBool[0:k+1])
		toReturn = append(toReturn, string(BoolArrayToBinaryString(prefix)))
	}
	return toReturn
}
