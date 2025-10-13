package blocks

import (
	"testing"
	"tp/common/keys"
)

func TestBlock(t *testing.T) {
	keyBlock := keys.GetKey("bloque1")
	keyEnd := keys.GetNullKey()
	block := GenerateBlockToStore([]byte{}, keyBlock, keyEnd)
	if !IsFinalBlock(block) {
		t.Errorf("Must be end")
	}
}
