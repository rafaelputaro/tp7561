package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tp/common"
)

const MSG_SERVER_SIGINT_ARRIVED = "SIGNIT arrived. Stopping Metrics Server"

func main() {
	config := LoadMetricsServerConfig()
	server := NewMetricsServer(config)
	// Iniciar el servidor
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		server.Serve()
		wg.Done()
	}()
	// Detener el servidor cuando llega la señal SIGINT
	handleSigintSignal(server)
	wg.Wait()
}

// Manejo de señal SIGINT
func handleSigintSignal(server *MetricsServer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		common.Log.Infof(MSG_SERVER_SIGINT_ARRIVED)
		server.DisposeMetricsServer()
	}()
}
