package dht

import (
	"errors"
	"fmt"
	"tp/common"
	"tp/common/keys"
	"tp/peer/dht/tiered_contact_storage"
	"tp/peer/helpers/file_manager"
)

const MSG_CONTACTS_ADDED_FOR_SEARCH = "contacts added for search: %v"

// Representa el resultado de procesar un contacto en la búsqueda de una clave
type processNextContactReturn struct {
	blockHasBeenFound bool
	endFileFound      bool
	nextBlockFound    []byte
	fileNameFound     string
	err               error
}

// Retorna una neuva nstancia de processNextContactReturn
func newProcessNextContactReturn(blockHasBeenFound bool, endFileFound bool, nextBlockFound []byte, fileNameFound string, err error) *processNextContactReturn {
	return &processNextContactReturn{
		blockHasBeenFound: blockHasBeenFound,
		endFileFound:      endFileFound,
		nextBlockFound:    nextBlockFound,
		fileNameFound:     fileNameFound,
		err:               err,
	}
}

// Desencola un contacto de contactStorage e intenta consultar acerca de la clave al mismo
// Coloca en el canal el resultado bajo el siguiente formato processNextContactReturn
func processNextContact(node *Node, key []byte, fileName string, contactStorage *tiered_contact_storage.TieredContactStorage, resChan chan processNextContactReturn, threadId int) {
	// Tomar el contacto más cercano y solicitar bloque/contactos cercanos
	contact, _ := contactStorage.Pop()
	// Si no hay más contactos retornar error
	if contact == nil {
		msg := fmt.Sprintf(MSG_ERROR_FILE_NOT_FOUND, fileName)
		resChan <- *newProcessNextContactReturn(false, false, keys.GetNullKey(), "", errors.New(msg))
		return
	}
	fileNameFound, nextBlockKeyFound, data, neighborContacts, err := node.SndFindBlock(node.Config, *contact, key)
	// Si hay un error retorna que no se encontró el bloque
	if err != nil {
		resChan <- *newProcessNextContactReturn(false, false, keys.GetNullKey(), fileNameFound, nil)
		return
	}
	// Agrego contacto a lista local
	node.scheduleAddContactTask(*contact)
	if len(fileNameFound) > 0 {
		endFile, err := file_manager.StoreBlockOnDownload(fileNameFound, data)
		if err == nil {
			common.Log.Debugf(MSG_FILE_FOUND, fileNameFound)
			resChan <- *newProcessNextContactReturn(true, endFile, nextBlockKeyFound, fileNameFound, nil)
		} else {
			resChan <- *newProcessNextContactReturn(false, false, keys.GetNullKey(), "", err)
		}
		return
	}
	// Agregar contactos encontrados a la búsqueda
	if neighborContacts != nil {
		common.Log.Debugf(MSG_CONTACTS_ADDED_FOR_SEARCH, len(neighborContacts))
		contactStorage.PushContacts(neighborContacts)
	}
	resChan <- *newProcessNextContactReturn(false, false, keys.GetNullKey(), "", nil)
}
