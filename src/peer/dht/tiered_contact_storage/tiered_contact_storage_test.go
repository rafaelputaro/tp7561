package tiered_contact_storage

import (
	"strconv"
	"testing"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
)

/*
	func TestStorageEasy(t *testing.T) {
		storage := NewTieredContactStorage(helpers.GetKey(""))
		if !storage.IsEmpty() {
			t.Errorf("Must be empty")
		}
		contact := contacts_queue.NewContact(helpers.GetKey("1"), "contacto1")
		storage.Push(*contact)
		if storage.IsEmpty() {
			t.Errorf("Must not be empty")
		}

		contactRec, tier := storage.Pop()
		if !storage.IsEmpty() {
			t.Errorf("Must be empty")
		}

		if contactRec == nil {
			t.Errorf("Contact must not be nill")
		} else {
			if contactRec.Url != contact.Url {
				t.Errorf("Incorrect contact: Found: %v | Expect: %v", contactRec.Url, contact.Url)
			}
		}
		if tier == 0 {
			t.Errorf("Incorrect tier: %v", tier)
		}
		contactRec, tier = storage.Pop()
		if !storage.IsEmpty() {
			t.Errorf("Must be empty")
		}
		if contactRec != nil {
			t.Errorf("Expect nil contact")
		}
		if tier != INVALID_TIER {
			t.Errorf("Expect tier: %v | Found: %v", INVALID_TIER, tier)
		}
	}
*/
func TestStorageHard(t *testing.T) {
	//rand.Seed(time.Now().UnixNano())
	storage := NewTieredContactStorage(helpers.GetKey(""))
	maxContacts := 10 //1000
	for contactNum := range maxContacts {
		contactNumStr := strconv.Itoa(contactNum)
		keyContact := helpers.GetKey(contactNumStr)
		contactStr := "contact" + contactNumStr
		println(helpers.GetLogDistance(helpers.GetKey(contactStr), helpers.GetKey("")))
		storage.Push(*contacts_queue.NewContact(keyContact, contactStr))
	}
	prevTier := 0
	tierFound := 0
	for countContacts := maxContacts; countContacts > 0; countContacts-- {
		if storage.count != countContacts {
			println("Expected: %v | Found: %v", countContacts, storage.count)
			//t.Errorf("Expected: %v | Found: %v", countContacts, storage.count)
		}
		_, tierFound = storage.Pop()
		/*if prevTier > tierFound {
			println("Out of order: Prev: " + strconv.Itoa(prevTier) + "| New: " + strconv.Itoa(tierFound) + "| Count: " + strconv.Itoa(countContacts))
			//t.Errorf("Out of order: Prev: %v | New: %v | Count: %v", prevTier, tierFound, countContacts)
		}
		*/
		//else {
		println("Prev: "+strconv.Itoa(prevTier)+" New: "+strconv.Itoa(tierFound)+" Count: ", countContacts)
		//}
		prevTier = tierFound

	}
	if !storage.IsEmpty() {
		t.Errorf("Must be empty")
	}
}
