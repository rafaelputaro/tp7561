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
const MSG_FAIL_ON_SHARE_CONTACTS = "error on sharing contacts: %v"
const MSG_FAIL_ON_SEND_STORE = "error sending store message: %v"
const MSG_FAIL_ON_SEND_FIND_BLOCK = "error sending find block message: %v"
const MSG_PING_ATTEMPT = "ping attempt: %v | error: %v"
const MSG_SHARE_CONTACTS_ATTEMPT = "share contacts attempt: %v | error: %v"
const MSG_STORE_ATTEMPT = "store block attempt: %v | error: %v"
const MSG_FIND_BLOCK_ATTEMPT = "find block attempt: %v | error: %v"
const MAX_RETRIES_ON_PING = 20
const MAX_RETRIES_ON_SHARE_CONTACTS_RECIP = 20
const MAX_RETRIES_ON_STORE = 20
const MAX_RETRIES_ON_FIND_BLOCK = 20

// Ping con retry. En caso de no poder efectuar el ping retorna error
type PingOp func(config helpers.PeerConfig, contact contacts_queue.Contact) error

// Share contact con retry. Retorna  <contacts><error>.En caso de no poder enviar el mensaje retorna error
type SndShareContactsRecipOp func(config helpers.PeerConfig, destContact contacts_queue.Contact, contacts []contacts_queue.Contact) ([]contacts_queue.Contact, error)

// Envío de store a un contacto con reintentos. Retorna <error>
type StoreOp func(config helpers.PeerConfig, contact contacts_queue.Contact, key []byte, value string, data []byte) error

// Find block con retry. Retorna <fileName><nextBlockKey><data con header><contacts><error> . En caso de no poder enviar el mensaje retorna error
type FindBlockOp func(config helpers.PeerConfig, destContact contacts_queue.Contact, key []byte) (string, []byte, []byte, []contacts_queue.Contact, error)

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

// Share contact con retry. Retorna  <contacts><error>.En caso de no poder enviar el mensaje retorna error
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
			var response *protopb.ShCtsRecipRes
			// compartir contacto
			response, err = client.ShCtsReciprocally(ctx, shContacOp)
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

// Envío de store a un contacto con reintentos. Retorna <error>
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

// Find block con retry. Retorna <fileName><nextBlockKey><data con header><contacts><error> . En caso de no poder enviar el mensaje retorna error
func SndFindBlock(config helpers.PeerConfig, destContact contacts_queue.Contact, key []byte) (string, []byte, []byte, []contacts_queue.Contact, error) {
	// conexión
	conn, client, ctx, cancel, err := communication.ConnectAsClient(destContact.Url, communication.LogFatalOnFailConnect)
	if err == nil {
		defer conn.Close()
		defer cancel()
		// armo los argumentos
		operands := protoUtils.CreateFindBlockOperands(config.Id, config.Url, key)
		// find block con retry
		for retry := range MAX_RETRIES_ON_FIND_BLOCK {
			var response *protopb.FindBlockRes
			// enviar mensaje
			response, err = client.FindBlock(ctx, operands)
			if err != nil {
				common.Log.Infof(MSG_FIND_BLOCK_ATTEMPT, retry, err)
				// esperar
				helpers.SleepBetweenRetriesShort()
				continue
			}
			fileName, nextBlockKey, data, contacts := protoUtils.ParseFindBlockResults(response)
			return fileName, nextBlockKey, data, contacts, nil
		}
		return "", nil, nil, nil, err
	}
	common.Log.Errorf(MSG_FAIL_ON_SEND_FIND_BLOCK, err)
	return "", nil, nil, nil, err
}
