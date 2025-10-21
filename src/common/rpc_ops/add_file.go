package rpc_ops_common

import (
	"tp/common"
	"tp/common/communication"
	filetransfer "tp/common/files_common/file_transfer"
	"tp/common/keys"
	"tp/protobuf/protoUtils"
)

// Envío de add a un contacto con reintentos. Retorna <key><error>
func AddFile(url string, fileName string, path string) ([]byte, error) {
	// conexión
	conn, client, ctx, cancel, err := communication.ConnectAsClientGRPC(url, communication.LogFatalOnFailConnectGRPC)
	key := keys.GetNullKey()
	urlDest := ""
	if err == nil {
		defer conn.Close()
		defer cancel()
		// armo los argumentos
		operands := protoUtils.CreateAddFileOperands(fileName)
		for retry := range MAX_RETRIES_ON_ADD_FILE {
			// enviar add file message
			response, errAf := client.AddFile(ctx, operands)
			if errAf != nil {
				common.Log.Infof(MSG_ADD_FILE_ATTEMPT, fileName, retry, errAf)
				// esperar
				common.SleepBetweenRetries()
				continue
			}
			key, urlDest = protoUtils.ParseAddFileResults(response)
			break
		}
		// si la urlDest es nula significa que el archivo no puede ser subido
		if urlDest == "" {
			common.Log.Debugf(MSG_ERROR_FILE_EXIST, fileName)
			return key, err
		}
		if err = filetransfer.SendFile(urlDest, fileName, path); err == nil {
			return key, err
		}
	}
	common.Log.Errorf(MSG_FAIL_ON_SEND_ADD_FILE, err)
	return key, err
}
