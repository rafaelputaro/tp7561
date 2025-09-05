package dht

import (
	"errors"
	"tp/peer/helpers"
)

const MSG_ERROR_ON_GET_VALUE = "error on get value from key value table"
const MSG_ERROR_ON_ADD_VALUE = "error on add key value"
const MSG_ERROR_ON_UPDATE_VALUE = "error on update value from key value table"
const EMPTY_VALUE = ""

// Es una table que contiene pares clave valor
type KeyValueTable struct {
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
func (table *KeyValueTable) Add(key []byte, value string) error {
	keyS := helpers.KeyToString(key)
	_, exists := table.Entries[keyS]
	if exists {
		return errors.New(MSG_ERROR_ON_ADD_VALUE)
	}
	table.Entries[keyS] = value
	return nil
}

// Remueve una clave de la tabla
func (table *KeyValueTable) Remove(key []byte) {
	delete(table.Entries, helpers.KeyToString(key))
}

// Obtiene el valor para una clave. En caso de no disponer la clave retorna error
func (table *KeyValueTable) GetValue(key []byte) (string, error) {
	if value, ok := table.Entries[helpers.KeyToString(key)]; ok {
		return value, nil
	}
	helpers.Log.Errorf(MSG_ERROR_ON_GET_VALUE)
	return EMPTY_VALUE, errors.New(MSG_ERROR_ON_GET_VALUE)
}

// Actualiza el valor para una clave. En caso de no disponer la clave retorna error
func (table *KeyValueTable) UpdateValue(key []byte, newValue string) error {
	if _, ok := table.Entries[helpers.KeyToString(key)]; ok {
		table.Entries[helpers.KeyToString(key)] = newValue
		return nil
	}
	helpers.Log.Errorf(MSG_ERROR_ON_UPDATE_VALUE)
	return errors.New(MSG_ERROR_ON_UPDATE_VALUE)
}

// Retorna verdadero si la table contiene cierta clave
func (table *KeyValueTable) HasKey(key []byte) bool {
	_, err := table.GetValue(key)
	return err != nil
}
