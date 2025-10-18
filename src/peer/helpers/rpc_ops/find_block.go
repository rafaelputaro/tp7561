package rpc_ops

import (
	"tp/common"
	"tp/common/communication"
	"tp/common/contact"
	"tp/peer/helpers"
	proto_utils_peer "tp/peer/helpers/proto_utils"

	"tp/protobuf/protopb"
)

// Find block con retry. Retorna <fileName><nextBlockKey><data con header><contacts><error> . En caso de no poder enviar el mensaje retorna error
type FindBlockOp func(config helpers.PeerConfig, destContact contact.Contact, key []byte) (string, []byte, []byte, []contact.Contact, error)

// Find block con retry. Retorna <fileName><nextBlockKey><data con header><contacts><error> . En caso de no poder enviar el mensaje retorna error
func SndFindBlock(config helpers.PeerConfig, destContact contact.Contact, key []byte) (string, []byte, []byte, []contact.Contact, error) {
	// conexión
	conn, client, ctx, cancel, err := communication.ConnectAsClientGRPC(destContact.Url, communication.LogFatalOnFailConnectGRPC)
	if err == nil {
		defer conn.Close()
		defer cancel()
		// armo los argumentos
		operands := proto_utils_peer.CreateFindBlockOperands(config.Id, config.UrlGRPC, key)
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
			fileName, nextBlockKey, data, contacts := proto_utils_peer.ParseFindBlockResults(response)
			return fileName, nextBlockKey, data, contacts, nil
		}
		return "", nil, nil, nil, err
	}
	common.Log.Errorf(MSG_FAIL_ON_SEND_FIND_BLOCK, err)
	return "", nil, nil, nil, err
}
