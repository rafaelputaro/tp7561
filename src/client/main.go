package main

import (
	"tp/client/helpers"
	"tp/common"
	"tp/common/communication/url"
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
	key, err := rpc_ops_common.AddFile(urlPeer, "filec-1-1.txt", helpers.GenerateInputFilePath(*config, "filec-1-1.txt"))
	if err != nil {
		common.Log.Errorf("Error: %v %v", urlPeer, err)
	} else {
		common.Log.Debugf("Key: %v", keys.KeyToLogFormatString(key))
	}
	key, err = rpc_ops_common.AddFile(urlPeer, "filec-1-1.txt", helpers.GenerateInputFilePath(*config, "filec-1-1.txt"))
	if err != nil {
		common.Log.Errorf("Error: %v %v", urlPeer, err)
	} else {
		common.Log.Debugf("Key ReUpload: %v", keys.KeyToLogFormatString(key))
	}
	rpc_ops_common.GetFile(config.Url, urlPeer, key)

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
