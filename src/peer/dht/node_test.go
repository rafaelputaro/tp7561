package dht

import (
	"errors"
	"testing"

	"tp/common/contact"
	"tp/common/keys"
	"tp/peer/helpers"
)

func TestKeys(t *testing.T) {
	id := []byte{128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	keys.GenerateKeysFromOtherTrees(id)

}

func PingOpWithError(config helpers.PeerConfig, contact contact.Contact) error {
	return errors.New("Error")
}

func PingOpWithoutError(config helpers.PeerConfig, contact contact.Contact) error {
	return nil
}

func StoreOpWithError(config helpers.PeerConfig, contact contact.Contact, key []byte, value string) error {
	return errors.New("Error")
}

func StoreOpWithoutError(config helpers.PeerConfig, contact contact.Contact, key []byte, value string) error {
	return errors.New("Error")
}
