package main

const MESSAGE_START = "Starting client..."

func main() {
	client, err := NewClient()
	if err == nil {
		client.Start()
	}

	/*
		// iniciar cliente
		common.Log.Info(MESSAGE_START)
		common.InitLogger()
		client := NewClient()
		// servicio de métricas
		metrics := client_metrics.MetricsServiceInstance
		go func() {
			metrics.Serve()
		}()
		// esperar a que la mayoría de los pares se inicialicen intercambiando contactos
		common.SleepOnStart(10 * client.Config.NumberOfPairs)
		// agregar archivos en peer-1
		keysAdded := [][]byte{}
		// to check
		check := map[string][]byte{}
		urlPeer := url.GenerateURLPeer(1)
		files_common.OpOverDir(helpers.GenerateInputFilePath(client.Config, ""),
			func(fileName string) error {
				key, err := rpc_ops_common.AddFile(urlPeer, fileName, helpers.GenerateInputFilePath(client.Config, fileName))
				if err == nil {
					keyS := keys.KeyToLogFormatString(key)
					common.Log.Debugf("File added: %v | %v | %v", fileName, keyS, urlPeer)
					keysAdded = append(keysAdded, key)
					check[fileName] = key
					metrics.IncUploadedFileCount()
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
		urlPeer = url.GenerateURLPeer(client.Config.NumberOfPairs / 2)
		// solicitar archivos a el último peer
		for _, key := range keysAdded {
			_, _, errGet := rpc_ops_common.GetFile(client.Config.Url, urlPeer, key)
			if errGet != nil {
				common.Log.Debugf("Error on get file %v", errGet)
			}
			common.SleepShort(client.Config.NumberOfPairs)
		}
		// chequear si llegaron todos los archivos
		for range 100 {
			for file := range check {
				_, err := os.Stat(helpers.GenerateDownloadPath(client.Config, file))
				if err == nil {
					delete(check, file)
				} else {
					accepted, pending, errGet := rpc_ops_common.GetFile(client.Config.Url, urlPeer, check[file])
					if errGet != nil {
						common.Log.Debugf("Error on get file %v", errGet)
						continue
					}
					if accepted {
						common.Log.Debugf("Retry to get file %v", file)
						continue
					}
					if pending {
						common.Log.Debugf("Pending file %v", file)
					}
				}
			}
			if len(check) == 0 {
				common.Log.Infof("Pass")
				return
			}
			common.SleepBetweenRetries()
		}
		common.Log.Infof("Fail")
	*/
}
