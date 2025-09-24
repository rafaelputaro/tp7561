package utils

import (
	"testing"
)

func TestKeys(t *testing.T) {
	ret := PathExists("/tmp")
	print("Retorno:", ret)
}
