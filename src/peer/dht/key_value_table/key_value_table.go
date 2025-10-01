package key_value_table

import (
	"errors"
	"os"
	"sync"
	"tp/common"
	"tp/peer/helpers"
	"tp/peer/helpers/file_manager"
)

const MSG_ERROR_ON_GET_VALUE = "error on get value from key value table"
const MSG_ERROR_ON_ADD_VALUE = "error on add key value"
const MSG_ERROR_ON_UPDATE_VALUE = "error on update value from key value table"
const EMPTY_VALUE = ""

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
	keyS := helpers.KeyToString(key)
	_, exists := table.Entries[keyS]
	if exists {
		return errors.New(MSG_ERROR_ON_ADD_VALUE)
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
	return nil
}

// Remueve una clave de la tabla
func (table *KeyValueTable) Remove(key []byte) {
	// tomar lock
	table.mutex.Lock()
	defer table.mutex.Unlock()
	// operar
	delete(table.Entries, helpers.KeyToString(key))
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
	data, err := file_manager.GetBlock(fileName)
	return fileName, data, err
}

// Obtiene el valor para una clave. En caso de no disponer la clave retorna error
func (table *KeyValueTable) getFileName(key []byte) (string, error) {
	if value, ok := table.Entries[helpers.KeyToString(key)]; ok {
		return value, nil
	}
	common.Log.Errorf(MSG_ERROR_ON_GET_VALUE)
	return EMPTY_VALUE, errors.New(MSG_ERROR_ON_GET_VALUE)
}

// Actualiza el valor para una clave. En caso de no disponer la clave retorna error
func (table *KeyValueTable) UpdateValue(key []byte, newValue string) error {
	// tomar lock
	table.mutex.Lock()
	defer table.mutex.Unlock()
	// operar
	if _, ok := table.Entries[helpers.KeyToString(key)]; ok {
		table.Entries[helpers.KeyToString(key)] = newValue
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
