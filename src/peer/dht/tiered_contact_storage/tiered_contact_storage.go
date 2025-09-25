package tiered_contact_storage

import (
	"tp/peer/dht/bucket_table/contacts_queue"
)

const MSG_ERROR_INVALID_TIER = "invalid tier error"

const INVALID_TIER = -1

type TieredContactStorage struct {
	tiers   map[int]*ContactStack
	maxTier int
	count   int
}

func NewTieredContactStorage() *TieredContactStorage {
	return &TieredContactStorage{
		tiers:   map[int]*ContactStack{},
		maxTier: INVALID_TIER,
		count:   0,
	}
}

func (storage *TieredContactStorage) Push(contact contacts_queue.Contact, tier int) {
	contactStack, exists := storage.tiers[tier]
	if !exists {
		contactStack = NewContactStack()
		storage.tiers[tier] = contactStack
	}
	contactStack.Push(contact)
	storage.count++
	if storage.maxTier < tier {
		storage.maxTier = tier
	}
}

func (storage *TieredContactStorage) Pop() (*contacts_queue.Contact, int) {
	if storage.IsEmpty() {
		return nil, INVALID_TIER
	}
	contactToReturn := storage.tiers[storage.maxTier].Pop()
	tierToReturn := storage.maxTier
	storage.count--
	storage.checkTiers()
	return contactToReturn, tierToReturn
}

func (storage *TieredContactStorage) IsEmpty() bool {
	return storage.count == 0
}

func (storage *TieredContactStorage) checkTiers() {
	for tier := storage.maxTier; tier >= 0; tier-- {
		contactStack, exists := storage.tiers[tier]
		if exists {
			if !contactStack.IsEmpty() {
				break
			}
		}
		storage.maxTier--
	}
}
