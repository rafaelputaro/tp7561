package rpc_ops

import (
	"tp/common"
	"tp/common/communication"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
	"tp/protobuf/protoUtils"
	"tp/protobuf/protopb"
)

// Share contact con retry. Retorna  <contacts><error>.En caso de no poder enviar el mensaje retorna error
type SndShareContactsRecipOp func(config helpers.PeerConfig, destContact contacts_queue.Contact, contacts []contacts_queue.Contact) ([]contacts_queue.Contact, error)

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
				common.SleepBetweenRetries()
				continue
			}
			return protoUtils.ParseShareContactsReciprocallyResults(response), nil
		}
		return nil, err
	}
	common.Log.Errorf(MSG_FAIL_ON_SHARE_CONTACTS, err)
	return nil, err
}
