package main

import (
	"tp/client/helpers"
	"tp/common"
	"tp/common/communication/url"
	"tp/common/files_common"
	filetransfer "tp/common/files_common/file_transfer"
	"tp/common/keys"
	rpc_ops_common "tp/common/rpc_ops"
)

const MESSAGE_START = "Starting client..."

func main() {
	// iniciar cliente
	common.Log.Info(MESSAGE_START)
	common.InitLogger()
	config := helpers.LoadConfig()
	helpers.InitStore(*config)
	// esperar a que la mayoría de los pares se inicialicen intercambiando contactos
	common.SleepOnStart(config.NumberOfPairs)
	// agregar archivos en peer-1
	keysAdded := [][]byte{}
	urlPeer := url.GenerateURLPeer(1)
	files_common.OpOverDir(helpers.GenerateInputFilePath(*config, ""),
		func(fileName string) error {
			key, err := rpc_ops_common.AddFile(urlPeer, fileName, helpers.GenerateInputFilePath(*config, fileName))
			if err == nil {
				common.Log.Debugf("File added: %v | %v | %v", fileName, keys.KeyToLogFormatString(key), urlPeer)
				keysAdded = append(keysAdded, key)
			}
			common.SleepShort(config.NumberOfPairs)
			return err
		})
	// escuchar llegada de archivos
	filetransfer.NewReceiver(
		config.Url,
		func(fileName string) string {
			return helpers.GenerateDownloadPath(*config, fileName)
		},
		func([]byte, string) {},
	)
	common.SleepOnStart(config.NumberOfPairs)
	urlPeer = url.GenerateURLPeer(config.NumberOfPairs - 1)
	// solicitar archivos
	for _, key := range keysAdded {
		_, errGet := rpc_ops_common.GetFile(config.Url, urlPeer, key)
		if errGet != nil {
			common.Log.Debugf("Error get file %v", errGet)
			common.SleepShort(config.NumberOfPairs)
		}
	}
	for {
		common.Log.Debugf("Listen")
		common.SleepBetweenRetries()
		common.SleepBetweenRetries()
	}

}
