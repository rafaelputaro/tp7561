package rpc_ops

import (
	"tp/common"
	"tp/common/communication"
	"tp/common/contact"
	"tp/peer/helpers"
	"tp/protobuf/protoUtils"
	"tp/protobuf/protopb"
)

// Share contact con retry. Retorna  <contacts><error>.En caso de no poder enviar el mensaje retorna error
type SndShareContactsRecipOp func(config helpers.PeerConfig, destContact contact.Contact, contacts []contact.Contact) ([]contact.Contact, error)

// Share contact con retry. Retorna  <contacts><error>.En caso de no poder enviar el mensaje retorna error
func SndShareContactsRecip(config helpers.PeerConfig, destContact contact.Contact, contacts []contact.Contact) ([]contact.Contact, error) {
	// conexión
	conn, client, ctx, cancel, err := communication.ConnectAsClientGRPC(destContact.Url, communication.LogFatalOnFailConnectGRPC)
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
