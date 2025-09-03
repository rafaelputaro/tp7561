package bucket_table

import (
	"fmt"
	"testing"
	"tp/peer/helpers"
)

func TestBucketTable(t *testing.T) {
	//key := helpers.GetKey("")
	key := []byte{}
	key = append(key, 4)
	prefixes := helpers.GeneratePrefixes(key)
	fmt.Println(prefixes)
	//arrayBool := helpers.ConvertToBoolArray(key)
	//print("%v", fmt.Sprintf("%v", arrayBool))
	//print("%v", len(arrayBool))
	//	table := NewBucketTable(key)

}
