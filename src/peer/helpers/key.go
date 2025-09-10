package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const EMPTY_KEY = ""
const LENGTH_KEY_IN_BITS = 256
const MSG_ERROR_ON_PARSE = "error on parse"

// Obtiene una key SHA256 desde un string
func GetKey(data string) []byte {
	h := sha256.New()
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

// Retorna una clave construída en base a un prefijo la cuál consiste en agregar ceros hasta
// completar los bits menos significativos
func PrefixToKey(prefix string) []byte {
	// completar la cadena
	keyStr := prefix + strings.Repeat(`0`, LENGTH_KEY_IN_BITS-len(prefix))
	toReturn := []byte{}
	for i := 0; i < LENGTH_KEY_IN_BITS; i += 8 {
		println("Len %v | Index %v", len(keyStr), i)
		/*
			intValue, err := strconv.ParseInt(keyStr[i:i+8], 2, 8)
			if err != nil {
				Log.Fatalf(MSG_ERROR_ON_PARSE)
			}
			toReturn = append(toReturn, byte(intValue))*/
	}
	return toReturn
}

// Transforma todos los prefijos de una lista en claves
func PrefixesToCompleteToKeys(prefixes []string) [][]byte {
	toReturn := [][]byte{}
	for _, prefix := range prefixes {
		toReturn = append(toReturn, PrefixToKey(prefix))
	}
	return toReturn
}

// Genera una lista que contiene una clave de cada uno de los árboles a los cuales
// no pertence la clave
func GenerateKeysFromOtherTrees(key []byte) [][]byte {
	prefixesStr := GeneratePrefixesOtherTrees(key)
	/*	println("Values %v", prefixesStr)
		for _, str := range prefixesStr {
			fmt.Printf("Elemento: %v\n", len(str)) // %s imprime la string, \n añade una nueva línea
		}
	*/
	return PrefixesToCompleteToKeys(prefixesStr)
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
