package main

import (
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
	// obtener contactos para la tabla propia en el inicio
	peer.SndShareContactsToBootstrap()
	// servir a resto de pares
	peer.Serve()
}
