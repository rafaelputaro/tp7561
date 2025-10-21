package main

import (
	"tp/client/helpers"
	"tp/common"
	"tp/common/communication/url"
	filetransfer "tp/common/file_transfer"
	"tp/common/keys"
	rpc_ops_common "tp/common/rpc_ops"
)

const MESSAGE_START = "Starting client..."

func main() {
	common.Log.Info(MESSAGE_START)
	common.InitLogger()
	config := helpers.LoadConfig()
	helpers.InitStore(*config)
	urlPeer := url.GenerateURLPeer(1)
	common.Log.Errorf("Url Peer: %v ", urlPeer)
	common.SleepBetweenRetries()
	common.SleepBetweenRetries()
	common.SleepBetweenRetries()
	common.SleepBetweenRetries()
	key, err := rpc_ops_common.AddFile(urlPeer, "filec-1-1.txt", helpers.GenerateInputFilePath(*config, "filec-1-1.txt"))
	common.Log.Debugf("filce-1-1.txt %v ", keys.KeyToLogFormatString(keys.GetKey("filec-1-1.txt")))
	if err != nil {
		common.Log.Errorf("Add File Error: %v %v", urlPeer, err)
	} else {
		common.Log.Debugf("Add File Key: %v", keys.KeyToLogFormatString(key))
	}
	/*
		key, err = rpc_ops_common.AddFile(urlPeer, "filec-1-1.txt", helpers.GenerateInputFilePath(*config, "filec-1-1.txt"))
		if err != nil {
			common.Log.Errorf("Error: %v %v", urlPeer, err)
		} else {
			common.Log.Debugf("Key ReUpload: %v", keys.KeyToLogFormatString(key))
		}*/

	_, errRec := filetransfer.NewReceiver(
		config.Url,
		func(fileName string) string {
			return helpers.GenerateDownloadPath(*config, fileName)
		},
		func([]byte, string) {},
	)
	if errRec != nil {
		common.Log.Debugf("Cant start receiver %v", errRec)
	}
	common.SleepBetweenRetries()
	common.SleepBetweenRetries()
	_, errGet := rpc_ops_common.GetFile(config.Url, urlPeer, key)
	if errGet != nil {
		common.Log.Debugf("Error get file %v", errGet)
	}
	if errRec == nil {
		for {
			common.Log.Debugf("Listen")
			common.SleepBetweenRetries()
			common.SleepBetweenRetries()
		}
	}

	/*
		utils.InitStore()
		// crear par
		peer := NewPeer(*config)
		wg := new(sync.WaitGroup)
		wg.Add(1)
		numberOfPairs := peer.Config.NumberOfPairs
		// obtener contactos para la tabla propia en el inicio
		go func() {
			helpers.SleepOnStart(numberOfPairs)
			peer.SndShCtsToBootstrap()
			helpers.SleepOnStart(numberOfPairs)
			file_manager.UploadLocalFiles(func(fileName string) error {
				peer.DoAddFile(fileName)
				return nil
			})
			helpers.SleepShort(numberOfPairs)
			if peer.NodeDHT.IsBootstrapNode() {
				for fileNum := 1; fileNum < 15; fileNum++ {
					if peer.GetFile("file-"+strconv.Itoa(fileNum)+"-1.txt") != nil {
						common.Log.Debugf("No se encontro archivo: %v", fileNum)
					}
					helpers.SleepShort(numberOfPairs)
				}
			}
			wg.Done()
		}()
		// servir a resto de pares
		peer.Serve()
		wg.Wait()
		peer.DisposePeer()

	*/
}
