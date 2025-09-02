package dht

import (
	"errors"
	"tp/peer/helpers"
)

const MSG_ERROR_ON_GET_VALUE = "error on get value from key value table"
const MSG_ERROR_ON_UPDATE_VALUE = "error on update value from key value table"
const EMPTY_VALUE = ""

// Es una table que contiene pares clave valor
type TKeyValueTable struct {
	Entries map[string]string
}

// Retorna una nueva instancia de tabla clave valor
func NewKeyValueTable() *TKeyValueTable {
	table := TKeyValueTable{
		Entries: map[string]string{},
	}
	return &table
}

// Agrega una nueva clave a la tabla
func (table *TKeyValueTable) Add(key []byte, value string) {
	table.Entries[helpers.KeyToString(key)] = value
}

// Remueve una clave de la tabla
func (table *TKeyValueTable) Remove(key []byte) {
	delete(table.Entries, helpers.KeyToString(key))
}

// Obtiene el valor para una clave. En caso de no disponer la clave retorna error
func (table *TKeyValueTable) GetValue(key []byte) (string, error) {
	if value, ok := table.Entries[helpers.KeyToString(key)]; ok {
		return value, nil
	}
	helpers.Log.Errorf(MSG_ERROR_ON_GET_VALUE)
	return EMPTY_VALUE, errors.New(MSG_ERROR_ON_GET_VALUE)
}

// Actualiza el valor para una clave. En caso de no disponer la clave retorna error
func (table *TKeyValueTable) UpdateValue(key []byte, newValue string) error {
	if _, ok := table.Entries[helpers.KeyToString(key)]; ok {
		table.Entries[helpers.KeyToString(key)] = newValue
		return nil
	}
	helpers.Log.Errorf(MSG_ERROR_ON_UPDATE_VALUE)
	return errors.New(MSG_ERROR_ON_UPDATE_VALUE)
}

// Retorna verdadero si la table contiene cierta clave
func (table *TKeyValueTable) HasKey(key []byte) bool {
	_, err := table.GetValue(key)
	return err != nil
}
