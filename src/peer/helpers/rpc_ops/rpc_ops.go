package rpc_ops

import (
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
)

type PingOp func(config helpers.PeerConfig, contact contacts_queue.Contact) error

type FindNodeOp func(config helpers.PeerConfig, contact contacts_queue.Contact, targetId []byte) ([]contacts_queue.Contact, error)

type FindValueOp func(config helpers.PeerConfig, contact contacts_queue.Contact, targetKey []byte) (string, error)

type StoreOp func(config helpers.PeerConfig, contact contacts_queue.Contact, key []byte, value string) error
