package tiered_contact_storage

import (
	"strconv"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
)

const MSG_ERROR_INVALID_TIER = "invalid tier error"

const INVALID_TIER = -1

// Representa un contendor de contactos los cuales se agrupan de acuerdo a la distancia a
// una clave dada. La extacción de contacto se hace iniciando por las menores distancias
type TieredContactStorage struct {
	tiers   map[int]*ContactStack
	minTier int
	count   int
	key     []byte
}

// Rertorna una instancia vacía lista para ser utilizada
func NewTieredContactStorage(key []byte) *TieredContactStorage {
	return &TieredContactStorage{
		tiers:   map[int]*ContactStack{},
		minTier: helpers.INFINITY_DISTANCE,
		count:   0,
		key:     key,
	}
}

// Inserta un nuevo contacto en base a la distancia respecto a la clave inicial
func (storage *TieredContactStorage) Push(contact contacts_queue.Contact) {
	tier, _ := helpers.GetLogDistance(contact.ID, storage.key)
	_, exists := storage.tiers[tier]
	if !exists {
		storage.tiers[tier] = NewContactStack()
	}
	storage.tiers[tier].Push(contact)
	storage.count++
	if storage.minTier > tier {
		storage.minTier = tier
	}
}

func (storage *TieredContactStorage) Pop() (*contacts_queue.Contact, int) {
	if storage.IsEmpty() {
		return nil, INVALID_TIER
	}
	contactToReturn := storage.tiers[storage.minTier].Pop()
	tierToReturn := storage.minTier
	storage.count--
	storage.checkTiers()
	return contactToReturn, tierToReturn
}

func (storage *TieredContactStorage) IsEmpty() bool {
	return storage.count == 0
}

func (storage *TieredContactStorage) checkTiers() {
	newMinTier := helpers.INFINITY_DISTANCE
	for tier := range storage.tiers {
		if storage.tiers[tier].IsEmpty() {
			delete(storage.tiers, tier)
			continue
		}
		newMinTier = tier
		break
	}
	storage.minTier = newMinTier
	println("Min tier: " + strconv.Itoa(storage.minTier) + " Count: " + strconv.Itoa(storage.count))
}
