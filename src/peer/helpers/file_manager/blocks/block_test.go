package blocks

import (
	"testing"
	"tp/peer/helpers"
)

func TestBlock(t *testing.T) {
	keyBlock := helpers.GetKey("bloque1")
	keyEnd := helpers.GetNullKey()
	block := GenerateBlockToStore([]byte{}, keyBlock, keyEnd)
	if !IsFinalBlock(block) {
		t.Errorf("Must be end")
	}
}
