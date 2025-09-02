package helpers

import (
	"crypto/sha1"
)

// Obtiene una key SHA1 desde un string
func GetKey(data string) []byte {
	h := sha1.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}

func KeyToString(key []byte) string {
	return string(key)
}
