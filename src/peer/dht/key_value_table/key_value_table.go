package key_value_table

import (
	"errors"
	"os"
	"sync"
	"tp/common"
	"tp/common/keys"
	"tp/peer/helpers/file_manager"
)

const MSG_ERROR_ON_GET_VALUE = "The file associated with the key was not found in the bucket table"
const MSG_ERROR_ON_ADD_VALUE = "error on add key value"
const MSG_ERROR_ON_UPDATE_VALUE = "error on update value from key value table"
const EMPTY_VALUE = ""
const MSG_LOG_KEY_VALUE = "key: %v | value: %v"
const MSG_KEY_ADD = "add key: %v | value: %v"

// Es una table que contiene pares clave valor y permite almacenar localmente bloques
// asociados a las claves como archivos donde los valores almacenados son los nombres
// de los archivos
type KeyValueTable struct {
	mutex   sync.Mutex
	Entries map[string]string
}

// Retorna una nueva instancia de tabla clave valor
func NewKeyValueTable() *KeyValueTable {
	table := KeyValueTable{
		Entries: map[string]string{},
	}
	return &table
}

// Agrega una nueva clave a la tabla. En caso de que la clave ya se encontraba
// en la tabla retorna error
func (table *KeyValueTable) Add(key []byte, fileName string, data []byte) error {
	// tomar lock
	table.mutex.Lock()
	defer table.mutex.Unlock()
	// verificar si existe la clave localmente
	keyS := keys.KeyToLogFormatString(key)
	_, exists := table.Entries[keyS]
	if exists {
		return errors.New(MSG_ERROR_ON_ADD_VALUE + ":" + fileName)
	}
	// guardar en disco
	err := file_manager.StoreBlock(fileName, data)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			common.Log.Debugf(MSG_ERROR_ON_ADD_VALUE + ":" + err.Error())
			return err
		}
	}
	// almacenar clave en tabla
	table.Entries[keyS] = fileName
	common.Log.Debugf(MSG_KEY_ADD, keyS, fileName)
	return nil
}

// Remueve una clave de la tabla
func (table *KeyValueTable) Remove(key []byte) {
	// tomar lock
	table.mutex.Lock()
	defer table.mutex.Unlock()
	// operar
	delete(table.Entries, keys.KeyToLogFormatString(key))
}

// Obtiene el nombre del archivo junto a sus datos para cierta clave. En caso de no disponer
// la clave retorna error. Los datos del archivo se retornan con su header completo.
// <fileName><data><error>
func (table *KeyValueTable) Get(key []byte) (string, []byte, error) {
	// obtener nombre del archivo
	fileName, err := table.getFileName(key)
	// si no lo contiene retornar err
	if err != nil {
		return "", nil, err
	}
	// leer datos del archivo
	data, err := file_manager.GetBlockFromStore(fileName)
	return fileName, data, err
}

// Verifica si un archivo existe en la KeyValueKey local
func (table *KeyValueTable) FileExistLocally(fileName string) bool {
	key := keys.GetKey(fileName)
	fileNameFound, err := table.getFileName(key)
	if err != nil {
		return false
	}
	return fileNameFound == fileName
}

// Obtiene el valor para una clave. En caso de no disponer la clave retorna error
func (table *KeyValueTable) getFileName(key []byte) (string, error) {
	if value, ok := table.Entries[keys.KeyToLogFormatString(key)]; ok {
		return value, nil
	}
	common.Log.Infof(MSG_ERROR_ON_GET_VALUE)
	return EMPTY_VALUE, errors.New(MSG_ERROR_ON_GET_VALUE)
}

// Actualiza el valor para una clave. En caso de no disponer la clave retorna error
func (table *KeyValueTable) UpdateValue(key []byte, newValue string) error {
	// tomar lock
	table.mutex.Lock()
	defer table.mutex.Unlock()
	// operar
	if _, ok := table.Entries[keys.KeyToLogFormatString(key)]; ok {
		table.Entries[keys.KeyToLogFormatString(key)] = newValue
		return nil
	}
	common.Log.Errorf(MSG_ERROR_ON_UPDATE_VALUE)
	return errors.New(MSG_ERROR_ON_UPDATE_VALUE)
}

// Retorna verdadero si la table contiene cierta clave
func (table *KeyValueTable) HasKey(key []byte) bool {
	_, err := table.getFileName(key)
	return err != nil
}

// Se loguean las claves contenidas en la tabla junto a los nombres de sus archivos
// asociados
func (table *KeyValueTable) LogKeysAndValues() {
	for key, value := range table.Entries {
		common.Log.Debugf(MSG_LOG_KEY_VALUE, key, value)
	}
}
