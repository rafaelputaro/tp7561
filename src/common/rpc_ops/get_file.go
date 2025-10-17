package rpc_ops_common

import (
	"tp/common"
	"tp/common/communication"
	"tp/protobuf/protoUtils"
)

// Envío de add a un contacto con reintentos. Retorna <key><error>
func GetFile(selfUrl string, destUrl string, key []byte) error {
	// conexión
	conn, client, ctx, cancel, err := communication.ConnectAsClientGRPC(destUrl, communication.LogFatalOnFailConnectGRPC)
	if err == nil {
		defer conn.Close()
		defer cancel()
		// envío con reintentos
		for retry := range MAX_RETRIES_ON_ADD_FILE {
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
			accepted, _ := protoUtils.ParseGetFileResults(response)
			common.Log.Debugf("accepted %v", accepted)
			break
		}

	}
	// @TODO continuar con esto
	common.Log.Errorf(MSG_FAIL_ON_SEND_GET_FILE, err)
	return err
}
