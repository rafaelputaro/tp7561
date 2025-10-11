package rpc_ops

import (
	"tp/common"
	"tp/common/communication"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
	"tp/protobuf/protoUtils"
)

// Envío de store a un contacto con reintentos. Retorna <error>
type StoreOp func(config helpers.PeerConfig, contact contacts_queue.Contact, key []byte, value string, data []byte) error

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
				common.SleepBetweenRetries()
				continue
			}
			return nil
		}
		return err
	}
	common.Log.Errorf(MSG_FAIL_ON_SEND_STORE, err)
	return err
}
