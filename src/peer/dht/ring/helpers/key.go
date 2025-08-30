package helpers

import (
	"crypto/sha1"
	"encoding/hex"
)

// / Obtiene una key SHA1 desde un string
func GetKey(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
