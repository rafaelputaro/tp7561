package main

import (
	"sync"
	"tp/common"
	"tp/peer/helpers"
	"tp/peer/helpers/file_manager/utils"
)

const MESSAGE_START = "Starting node..."

func main() {
	common.Log.Info(MESSAGE_START)
	common.InitLogger()
	config := helpers.LoadConfig()
	utils.InitStore()
	// crear par
	peer := NewPeer(*config)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	numberOfPairs := peer.Config.NumberOfPairs
	// obtener contactos para la tabla propia en el inicio
	go func() {
		common.SleepOnStart(numberOfPairs)
		peer.SndShCtsToBootstrap()
		common.SleepOnStart(numberOfPairs)
		wg.Done()
	}()
	// servir a resto de pares
	peer.Serve()
	wg.Wait()
	peer.DisposePeer()
}
