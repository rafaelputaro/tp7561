package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

const EMPTY_KEY = ""
const LENGTH_KEY_IN_BITS = 256
const LENGTH_IN_BYTES = LENGTH_KEY_IN_BITS / 8
const MSG_ERROR_ON_PARSE = "error on parse"

// Obtiene una key SHA256 desde un string
func GetKey(data string) []byte {
	h := sha256.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}

// Transforma una clave dada como una cadena de bytes en un string
func KeyToString(key []byte) string {
	return string(key)
}

// Transforma una clave dada como una cadena de bytes en un string dado por su representación en
// hexadecimal
func KeyToHexString(key []byte) string {
	hexString := hex.EncodeToString(key)
	return hexString
}

// Transforma una clave en su representación en bits como un string
func ConvertToBinaryString(key []byte) string {
	boolArray := ConvertToBoolArray(key)
	return BoolArrayToBinaryString(boolArray)
}

// Transforma una clave en un array de booleanos
func ConvertToBoolArray(data []byte) []bool {
	res := make([]bool, len(data)*8)
	for i := range res {
		res[i] = data[i/8]&(0x80>>byte(i&0x7)) != 0
	}
	return res
}

// Genera una lista que contiene una clave de cada uno de los árboles a los cuales
// no pertenece la clave
func GenerateKeysFromOtherTrees(key []byte) [][]byte {
	toReturn := [][]byte{}
	// para cada byte
	for nByte, aByte := range key {
		// cambiar cada uno de los bits
		for i := range 8 {
			var operand byte = 128 >> i
			byteResultXor := aByte ^ operand
			newKey := []byte{}
			// completar bytes bajos
			for range nByte {
				newKey = append(newKey, 0)
			}
			// agregar byte calculado
			newKey = append(newKey, byteResultXor)
			// completar bytes altos
			for range LENGTH_IN_BYTES - nByte - 1 {
				newKey = append(newKey, 0)
			}
			// agregar nueva clave
			toReturn = append(toReturn, newKey)
		}
	}
	return toReturn
}

// Transforma la clave en un string en el formato dado por la configuración de inicio del módulo
func KeyToLogFormatString(key []byte) string {
	switch LoginFormatForKeys {
	case string(HEXA):
		return KeyToHexString(key)
	default:
		return ConvertToBinaryString(key)
	}
}

// Convierte un array de bool en una cadena con su representación binaria
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
func GeneratePrefixesOtherTreesAsStrings(key []byte) []string {
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
