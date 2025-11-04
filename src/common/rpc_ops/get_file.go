package rpc_ops_common

import (
	"fmt"
	"tp/common"
	"tp/common/communication"
	"tp/common/keys"
	"tp/protobuf/protoUtils"
)

const MSG_GET_FILE_ACCEPTED = "get file accepted: key %v | selfUrl %v | destUrl %v"

// Envío de add a un contacto con reintentos. Retorna <accepted><pending><error>
func GetFile(selfUrl string, destUrl string, key []byte) (bool, bool, error) {
	// conexión con reintentos
	conn, client, ctx, cancel, err := communication.ConnectAsClientGRPC(destUrl, communication.LogFatalOnFailConnectGRPC)
	if err == nil {
		defer conn.Close()
		defer cancel()
		// envío con reintentos
		for retry := range MAX_RETRIES_ON_GET_FILE {
			// armo los argumentos
			operands := protoUtils.CreateGetFileOperands(key, selfUrl)
			// enviar get file message
			response, errGf := client.GetFile(ctx, operands)
			if errGf != nil {
				common.Log.Infof(MSG_GET_FILE_ATTEMPT, key, retry, errGf)
				// esperar
				common.SleepBetweenRetries()
				continue
			}
			accepted, pending := protoUtils.ParseGetFileResults(response)
			if accepted {
				common.Log.Debugf(MSG_GET_FILE_ACCEPTED, keys.KeyToLogFormatString(key), selfUrl, destUrl)
			}
			return accepted, pending, errGf
		}
		common.Log.Errorf(MSG_FAIL_ON_SEND_GET_FILE, err)
		return false, false, fmt.Errorf(MSG_FAIL_ON_SEND_GET_FILE, keys.KeyToLogFormatString(key))
	}
	common.Log.Errorf(MSG_FAIL_ON_SEND_GET_FILE, err)
	return false, false, err
}
