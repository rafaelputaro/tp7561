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
	prefixes := helpers.GeneratePrefixesOtherTrees(key)
	fmt.Println(prefixes)
	//arrayBool := helpers.ConvertToBoolArray(key)
	//print("%v", fmt.Sprintf("%v", arrayBool))
	//print("%v", len(arrayBool))
	callback := func(url string) bool {
		return false
	}
	table := NewBucketTable(key, 10, callback)
	contact := []byte{}
	contact = append(contact, 5)
	table.AddContact(contact, "contact5::5051")
	contacts := table.GetContactsForId(contact)
	println("contact: %v", contacts[0].Url)
	/*
		contactF := table.DequeueContact("00000101")
		if contactF == nil {
			println("No hay contactos")
		} else {
			println("%v", contactF.ID)
		}*/

}
