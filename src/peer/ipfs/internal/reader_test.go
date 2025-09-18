package internal

import (
	"testing"
	"tp/common"
)

func TestBlockReader(t *testing.T) {
	fr, err := NewFileReader("/tmp/reader.txt")
	if err != nil {
		common.Log.Infof("error on open")
	}
	for {
		b, n, eof, err := fr.Next()
		if eof || err != nil {
			break
		}
		common.Log.Infof("byte: %v nbyte: %v eof: %v err: %v", string(b), n, eof, err)
	}

}
