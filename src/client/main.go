package main

import (
	"os"
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
	client := newClient()
	// servicio de métricas
	/*
		metrics := client_metrics.NewMetricsServer(func() int { return client.UploadedFileCount })

		go func() {
			metrics.Serve()
		}()
	*/
	// esperar a que la mayoría de los pares se inicialicen intercambiando contactos
	common.SleepOnStart(client.Config.NumberOfPairs)
	// agregar archivos en peer-1
	keysAdded := [][]byte{}
	// to check
	check := map[string]bool{}
	urlPeer := url.GenerateURLPeer(1)
	files_common.OpOverDir(helpers.GenerateInputFilePath(client.Config, ""),
		func(fileName string) error {
			key, err := rpc_ops_common.AddFile(urlPeer, fileName, helpers.GenerateInputFilePath(client.Config, fileName))
			if err == nil {
				keyS := keys.KeyToLogFormatString(key)
				common.Log.Debugf("File added: %v | %v | %v", fileName, keyS, urlPeer)
				keysAdded = append(keysAdded, key)
				check[fileName] = false
				client.UploadedFileCount++
			}
			common.SleepShort(client.Config.NumberOfPairs)
			return err
		})
	// escuchar llegada de archivos
	filetransfer.NewReceiver(
		client.Config.Url,
		func(fileName string) string {
			return helpers.GenerateDownloadPath(client.Config, fileName)
		},
		func([]byte, string) {},
	)
	common.SleepOnStart(client.Config.NumberOfPairs)
	urlPeer = url.GenerateURLPeer(client.Config.NumberOfPairs - 1)
	// solicitar archivos a el último peer
	for _, key := range keysAdded {
		_, errGet := rpc_ops_common.GetFile(client.Config.Url, urlPeer, key)
		if errGet != nil {
			common.Log.Debugf("Error get file %v", errGet)
		}
		common.SleepShort(client.Config.NumberOfPairs)
	}
	// chequear si llegaron todos los archivos
	for range 100 {
		for file := range check {
			_, err := os.Stat(helpers.GenerateDownloadPath(client.Config, file))
			if err == nil {
				delete(check, file)
			}
		}
		if len(check) == 0 {
			common.Log.Infof("Pass")
			return
		}
		common.SleepBetweenRetries()
	}
	common.Log.Infof("Fail")
}
