package tiered_contact_storage

import (
	"sort"
	"sync"
	"tp/common/keys"
	"tp/peer/dht/bucket_table/contacts_queue"
)

const MSG_ERROR_INVALID_TIER = "invalid tier error"

const INVALID_TIER = -1

// Representa un contendor de contactos los cuales se agrupan de acuerdo a la distancia a
// una clave dada. La extacción de contacto se hace iniciando por las menores distancias
type TieredContactStorage struct {
	contactTiers map[int]*ContactStack
	minTier      int
	count        int
	key          []byte
	presentTiers map[int]bool
	tierList     []int
	mutex        sync.Mutex
}

// Rertorna una instancia vacía lista para ser utilizada
func NewTieredContactStorage(key []byte) *TieredContactStorage {
	return &TieredContactStorage{
		contactTiers: map[int]*ContactStack{},
		minTier:      keys.INFINITY_DISTANCE,
		count:        0,
		key:          key,
		presentTiers: map[int]bool{},
		tierList:     []int{},
	}
}

// Inserta un nuevo contacto en base a la distancia respecto a la clave inicial
func (storage *TieredContactStorage) Push(contact contacts_queue.Contact) bool {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	return storage.doPush(contact)
}

// Inserta un nuevo contacto en base a la distancia respecto a la clave inicial
func (storage *TieredContactStorage) PushContacts(contacts []contacts_queue.Contact) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	for _, contact := range contacts {
		storage.doPush(contact)
	}
}

// Inserta un nuevo contacto en base a la distancia respecto a la clave inicial
func (storage *TieredContactStorage) doPush(contact contacts_queue.Contact) bool {
	tier, _ := keys.GetLogDistance(contact.ID, storage.key)
	_, exists := storage.contactTiers[tier]
	if !exists {
		storage.contactTiers[tier] = NewContactStack()
	}
	if storage.contactTiers[tier].Push(contact) {
		storage.tryAddToTierList(tier)
		storage.count++
		if storage.minTier > tier {
			storage.minTier = tier
		}
		return true
	}
	return false
}

// Intenta agregar un nuevo nivel asociada a una distancia
func (storage *TieredContactStorage) tryAddToTierList(tier int) {
	if _, ok := storage.presentTiers[tier]; !ok {
		storage.presentTiers[tier] = true
		storage.tierList = append(storage.tierList, tier)
	}
}

// Remueve el nivel asociado a la distancia más baja
func (storage *TieredContactStorage) removeBottomFromTierList() {
	if len(storage.tierList) > 0 {
		bottom := storage.tierList[0]
		delete(storage.contactTiers, bottom)
		delete(storage.presentTiers, bottom)
		storage.tierList = storage.tierList[1:]
	}
}

// Ordena la lista de distancias y la retorna
func (storage *TieredContactStorage) getSortedTiersList() []int {
	sort.Slice(storage.tierList, func(i, j int) bool {
		return storage.tierList[i] < storage.tierList[j]
	})
	return storage.tierList
}

// Obtener el contacto con la distancia más chica a la clave
func (storage *TieredContactStorage) Pop() (*contacts_queue.Contact, int) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	if storage.IsEmpty() {
		return nil, INVALID_TIER
	}
	contactToReturn := storage.contactTiers[storage.minTier].Pop()
	tierToReturn := storage.minTier
	storage.count--
	storage.checkTiers()
	return contactToReturn, tierToReturn
}

// Retorna verdadero si la instancia no contiene más contactos
func (storage *TieredContactStorage) IsEmpty() bool {
	return storage.count == 0
}

// Chequea que nivel tiene la distancia más chica actualizando el mínimo y removiendo
// niveles vacíos
func (storage *TieredContactStorage) checkTiers() {
	newMinTier := keys.INFINITY_DISTANCE
	// obtener lista de distancias ordenada
	tiers := storage.getSortedTiersList()
	for _, tier := range tiers {
		if storage.contactTiers[tier].IsEmpty() {
			storage.removeBottomFromTierList()
			continue
		}
		newMinTier = tier
		break
	}
	storage.minTier = newMinTier
}

// Retorna la cantidad de contactos
func (storage *TieredContactStorage) Count() int {
	return storage.count
}
