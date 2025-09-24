package main

import (
	"sync"
	"tp/common"
	"tp/peer/helpers"
	"tp/peer/helpers/file_manager"
)

const MESSAGE_START = "Starting node..."

func main() {
	common.Log.Info(MESSAGE_START)
	common.InitLogger()
	config := helpers.LoadConfig()
	file_manager.InitStore()
	// crear par
	peer := NewPeer(*config)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	// obtener contactos para la tabla propia en el inicio
	go func() {
		helpers.SleepOnStart()
		peer.SndShareContactsToBootstrap()
		helpers.SleepOnStart()
		helpers.SleepOnStart()
		if peer.NodeDHT.IsBootstrapNode() {
			peer.AddFile("file-1-1.txt")
		}
		helpers.SleepOnStart()
		if peer.NodeDHT.IsBootstrapNode() {
			peer.GetFile("file-1-1.txt")
		}
		wg.Done()
	}()
	// servir a resto de pares
	peer.Serve()
	wg.Wait()
}
