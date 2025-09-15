package main

import (
	"math/rand"
	"time"

	"tp/peer/helpers"
)

const MESSAGE_START = "Starting node..."

func main() {
	helpers.Log.Info(MESSAGE_START)
	helpers.InitLogger()
	config := helpers.LoadConfig()

	peer := NewPeer(*config)

	// Sleep porque aún no he agregado retry
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := randSource.Intn(30)
	t := time.Duration(r+15) * time.Second
	helpers.Log.Debugf("Tiempo: %v", t)
	time.Sleep(t)

	peer.SndShareContactsToBootstrap()

	peer.Serve()

	// @Todo utilizar Mutex sobre Peer. Crear una gofunc para servir mientras se arma la tabla
	// y se envían los archivos
}
