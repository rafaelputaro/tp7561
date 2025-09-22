package rpc_ops

import (
	"tp/common"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
	"tp/peer/helpers/communication"

	"tp/protobuf/protoUtils"
	"tp/protobuf/protopb"
)

const MSG_FAIL_ON_SEND_PING = "error sending ping: %v"
const MSG_PING_ATTEMPT = "ping attempt: %v | error: %v"
const MSG_SHARE_CONTACTS_ATTEMPT = "share contacts attempt: %v | error: %v"
const MSG_FAIL_ON_SHARE_CONTACTS = "error on sharing contacts: %v"
const MSG_FAIL_ON_SEND_STORE = "error sending store message: %v"
const MSG_STORE_ATTEMPT = "store block attempt: %v | error: %v"
const MAX_RETRIES_ON_PING = 20
const MAX_RETRIES_ON_SHARE_CONTACTS_RECIP = 20
const MAX_RETRIES_ON_STORE = 20

type PingOp func(config helpers.PeerConfig, contact contacts_queue.Contact) error

type SndShareContactsRecipOp func(config helpers.PeerConfig, destContact contacts_queue.Contact, contacts []contacts_queue.Contact) ([]contacts_queue.Contact, error)

type StoreOp func(config helpers.PeerConfig, contact contacts_queue.Contact, key []byte, value string, data []byte) error

// Ping con retry. En caso de no poder efectuar el ping retorna error
func SndPing(config helpers.PeerConfig, contact contacts_queue.Contact) error {
	// conexión
	conn, client, ctx, cancel, err := communication.ConnectAsClient(contact.Url, communication.LogFatalOnFailConnect)
	if err == nil {
		defer conn.Close()
		defer cancel()
		// ping con retry
		for retry := range MAX_RETRIES_ON_PING {
			_, err = client.Ping(ctx, protoUtils.CreatePingOperands(config.Id, config.Url))
			if err != nil {
				common.Log.Infof(MSG_PING_ATTEMPT, retry, err)
				// esperar
				helpers.SleepBetweenRetries()
				continue
			}
			return nil
		}
		return err
	}
	common.Log.Errorf(MSG_FAIL_ON_SEND_PING, err)
	return err
}

// Share contact con retry. En caso de no poder efectuar el ping retorna error
func SndShareContactsRecip(config helpers.PeerConfig, destContact contacts_queue.Contact, contacts []contacts_queue.Contact) ([]contacts_queue.Contact, error) {
	// conexión
	conn, client, ctx, cancel, err := communication.ConnectAsClient(destContact.Url, communication.LogFatalOnFailConnect)
	if err == nil {
		defer conn.Close()
		defer cancel()
		// share contact con retry
		for retry := range MAX_RETRIES_ON_SHARE_CONTACTS_RECIP {
			// armo los argumentos
			shContacOp := protoUtils.CreateShareContactsReciprocallyOperands(destContact, contacts)
			// compartir contacto
			var response *protopb.ShareContactsReciprocallyResults
			// compartir contacto
			response, err = client.ShareContactsReciprocally(ctx, shContacOp)
			if err != nil {
				common.Log.Infof(MSG_SHARE_CONTACTS_ATTEMPT, retry, err)
				// esperar
				helpers.SleepBetweenRetries()
				continue
			}
			return protoUtils.ParseShareContactsReciprocallyResults(response), nil
		}
		return nil, err
	}
	common.Log.Errorf(MSG_FAIL_ON_SHARE_CONTACTS, err)
	return nil, err
}

// Envío de store a un contacto con reintentos
func SndStore(config helpers.PeerConfig, contact contacts_queue.Contact, key []byte, blockName string, data []byte) error {
	// conexión
	conn, client, ctx, cancel, err := communication.ConnectAsClient(contact.Url, communication.LogFatalOnFailConnect)
	if err == nil {
		defer conn.Close()
		defer cancel()
		// send store
		for retry := range MAX_RETRIES_ON_STORE {
			// armo los argumentos
			operands := protoUtils.CreateStoreBlockOperands(config.Id, config.Url, key, blockName, data)
			// enviar store message
			_, err = client.StoreBlock(ctx, operands)
			if err != nil {
				common.Log.Infof(MSG_STORE_ATTEMPT, retry, err)
				// esperar
				helpers.SleepBetweenRetries()
				continue
			}
			return nil
		}
		return err
	}
	common.Log.Errorf(MSG_FAIL_ON_SEND_STORE, err)
	return err
}
