package main

import (
	"sync"
	"tp/common"
	"tp/peer/helpers"
)

const MESSAGE_START = "Starting node..."

func main() {
	common.Log.Info(MESSAGE_START)
	common.InitLogger()
	config := helpers.LoadConfig()
	// crear par
	peer := NewPeer(*config)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	// obtener contactos para la tabla propia en el inicio
	go func() {
		helpers.SleepOnStart()
		peer.SndShareContactsToBootstrap()
		wg.Done()
	}()
	// servir a resto de pares
	peer.Serve()
	wg.Wait()
}
