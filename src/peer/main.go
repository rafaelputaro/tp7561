package main

import (
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tp/peer/helpers"
	"tp/peer/protobuf/protopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const MESSAGE_START = "Starting node..."

func main() {
	helpers.Log.Info(MESSAGE_START)
	helpers.InitLogger()
	config := helpers.LoadConfig()

	peer := NewPeer(*config)

	lis, err := net.Listen("tcp", config.Url) //":"+config.Port)
	if err != nil {
		helpers.Log.Fatalf("failed to listen: %v", err)
	}
	// New server
	server := grpc.NewServer()
	protopb.RegisterOperationsServer(server, peer)
	// Register reflection service on gRPC server.
	reflection.Register(server)

	// Stop the gRPC server when the SIGINT signal arrives
	handleSigintSignal(server)

	helpers.Log.Infof("[SERVER] Starting gRPC Server")

	// Sleep porque aún no he agregado retry
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := randSource.Intn(30)
	t := time.Duration(r+15) * time.Second
	helpers.Log.Debugf("Tiempo: %v", t)
	time.Sleep(t)

	peer.SndShareContactsToBootstrap()

	if err := server.Serve(lis); err != nil {
		helpers.Log.Fatalf("failed to serve: %v", err)
	}
	helpers.Log.Infof("[SERVER] gRPC server has been stopped")

}

func handleSigintSignal(server *grpc.Server) {
	c := make(chan os.Signal, syscall.SIGINT)
	signal.Notify(c, os.Interrupt)
	go func() {
		// Block until we receive the SIGINT signal
		<-c
		helpers.Log.Infof("[SERVER] SIGNIT arrived. Stopping gRPC Server")
		server.Stop()
	}()
}
