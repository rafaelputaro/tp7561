package tiered_contact_storage

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
)

func TestStorageLow(t *testing.T) {
	storage := NewTieredContactStorage()
	if !storage.IsEmpty() {
		t.Errorf("Must be empty")
	}
	contact := contacts_queue.NewContact(helpers.GetKey("1"), "contacto1")
	storage.Push(*contact, 0)
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
	if tier != 0 {
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

func TestStorageHigh(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	storage := NewTieredContactStorage()
	maxContacts := 1000
	maxTiers := 30
	for contactNum := range maxContacts {
		tierNum := rand.Intn(maxTiers)
		contactNumStr := strconv.Itoa(contactNum)
		keyContact := helpers.GetKey(contactNumStr)
		storage.Push(*contacts_queue.NewContact(keyContact, "contact"+contactNumStr+" "+strconv.Itoa(tierNum)), tierNum)
	}
	lastTier := maxTiers
	for countContacts := maxContacts; countContacts > 0; countContacts-- {
		if storage.count != countContacts {
			t.Errorf("Expected: %v | Found: %v", countContacts, storage.count)
		}
		contact, tierFound := storage.Pop()
		if lastTier < tierFound {
			t.Errorf("Out of order: Last: %v | Found: %v", lastTier, tierFound)
		}
		lastTier = tierFound
		println(strconv.Itoa(lastTier) + " " + contact.Url)
	}
	if !storage.IsEmpty() {
		t.Errorf("Must be empty")
	}
}
