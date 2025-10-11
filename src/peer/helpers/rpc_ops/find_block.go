package rpc_ops

import (
	"tp/common"
	"tp/common/communication"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"

	"tp/protobuf/protoUtils"
	"tp/protobuf/protopb"
)

// Find block con retry. Retorna <fileName><nextBlockKey><data con header><contacts><error> . En caso de no poder enviar el mensaje retorna error
type FindBlockOp func(config helpers.PeerConfig, destContact contacts_queue.Contact, key []byte) (string, []byte, []byte, []contacts_queue.Contact, error)

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
				common.SleepBetweenRetriesShort()
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
