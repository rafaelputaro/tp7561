package rpc_ops_common

import (
	"tp/common"
	"tp/common/communication"
	"tp/common/files_common"
	"tp/common/files_common/uploader"
	"tp/common/keys"
	"tp/protobuf/protoUtils"
)

// Envío de add a un contacto con reintentos. Retorna <key><error>
func AddFile(url string, fileName string, path string) ([]byte, error) {
	// conexión
	conn, client, ctx, cancel, err := communication.ConnectAsClient(url, communication.LogFatalOnFailConnect)
	key := keys.GetNullKey()
	if err == nil {
		defer conn.Close()
		defer cancel()
		// crear un file reader
		reader, err := files_common.NewFileReader(path, uploader.MAX_PART_SIZE)
		if err != nil {
			return key, err
		}
		// leer y enviar file parts
		eof := false
		for !eof {
			// leer siguiente bloque
			dataR, partNum, eofR, errR := reader.Next()
			if errR != nil {
				return key, err
			}
			eof = eofR
			// envío con reintentos
			for retry := range MAX_RETRIES_ON_ADD_FILE {
				// armo los argumentos
				operands := protoUtils.CreateAddFileOperands(fileName, int32(partNum), dataR, eofR)
				// enviar add file message
				response, errAf := client.AddFile(ctx, operands)
				if errAf != nil {
					common.Log.Infof(MSG_ADD_FILE_ATTEMPT, fileName, partNum, retry, errAf)
					// esperar
					common.SleepBetweenRetries()
					continue
				}
				key = protoUtils.ParseAddFileResults(response)
				break
			}
		}
		return key, err
	}
	common.Log.Errorf(MSG_FAIL_ON_SEND_ADD_FILE, err)
	return key, err
}
