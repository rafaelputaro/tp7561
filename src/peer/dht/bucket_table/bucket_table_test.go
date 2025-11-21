package bucket_table

import (
	"fmt"
	"testing"
	"tp/common/keys"
)

func TestBucketTable(t *testing.T) {
	//key := helpers.GetKey("")
	key := []byte{}
	key = append(key, 4)
	prefixes := keys.GeneratePrefixesOtherTreesAsStrings(key)
	fmt.Println(prefixes)
}
