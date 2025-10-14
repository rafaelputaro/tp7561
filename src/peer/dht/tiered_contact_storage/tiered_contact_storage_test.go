package tiered_contact_storage

import (
	"math/rand"
	"strconv"
	"testing"
	"tp/common/contact"
	"tp/common/keys"
)

func TestStorageEasy(t *testing.T) {
	storage := NewTieredContactStorage(keys.GetKey(""))
	if !storage.IsEmpty() {
		t.Errorf("Must be empty")
	}
	contact := contact.NewContact(keys.GetKey("1"), "contacto1")
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

func TestStorageHard(t *testing.T) {
	storage := NewTieredContactStorage(keys.GetKey(""))
	maxContacts := 10000
	for i := range maxContacts {
		contactNumStr := strconv.Itoa(rand.Intn(5000))
		keyContact := keys.GetKey(contactNumStr)
		contactStr := "contact" + strconv.Itoa(i) + " " + contactNumStr
		if !storage.Push(*contact.NewContact(keyContact, contactStr)) {
			t.Errorf("The contact could not be added: %v", contactStr)
		}
		if storage.Push(*contact.NewContact(keyContact, contactStr)) {
			t.Errorf("The contact was added twice: %v", contactStr)
		}
	}
	prevTier := 0
	tierFound := 0
	for countContacts := maxContacts; countContacts > 0; countContacts-- {
		if storage.Count() != countContacts {
			t.Errorf("Expected: %v | Found: %v", countContacts, storage.Count())
		}
		_, tierFound = storage.Pop()
		if prevTier > tierFound {
			t.Errorf("Out of order: Prev: %v | New: %v | Count: %v", prevTier, tierFound, countContacts)
		}
		prevTier = tierFound
	}
	if !storage.IsEmpty() {
		t.Errorf("Must be empty")
	}
}
