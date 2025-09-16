package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"tp/common"
	"tp/peer/helpers"
	"tp/protobuf/protopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const MSG_FAILED_TO_LISTEN = "failed to listen: %v"
const MSG_SERVER_GRPC_STARTING = "Starting gRPC Server"
const MSG_SERVER_SIGINT_ARRIVED = "SIGNIT arrived. Stopping gRPC Server"
const MSG_SERVER_STOPPED = "gRPC server has been stopped"
const MSG_FAILED_TO_SERVE = "failed to serve: %v"
const MSG_RETRY_LISTEN = "Listening retry"
const MAX_RETRY_LISTEN = 10
const MAX_RETRY_SERVE = 10

// Implementa la funcionalidad de grpc server para el par
type PeerService struct {
	Listener net.Listener
	Server   *grpc.Server
}

// Retorna una nueva instancia de Peer Service lista para ser utilizada
func NewPeerService(peer *Peer) *PeerService {
	var lis net.Listener
	var err error = nil
	for range MAX_RETRY_LISTEN {
		lis, err = net.Listen("tcp", peer.Config.Url)
		if err == nil {
			break
		}
		helpers.SleepBetweenRetries()
		common.Log.Debugf(MSG_RETRY_LISTEN)
	}
	if err != nil {
		common.Log.Fatalf(MSG_FAILED_TO_LISTEN, err)
	}
	// Nuevo servicio
	server := grpc.NewServer()
	protopb.RegisterOperationsServer(server, peer)
	// Registrar el servicio de reflexión en el servidor gRPC.
	reflection.Register(server)
	// Detener el servidor gRPC cuando llega la señal SIGINT
	handleSigintSignal(server)
	// Servidor inicializado
	common.Log.Infof(MSG_SERVER_GRPC_STARTING)
	return &PeerService{
		Listener: lis,
		Server:   server,
	}
}

// Manejo de señal SIGINT
func handleSigintSignal(server *grpc.Server) {
	c := make(chan os.Signal, syscall.SIGINT)
	signal.Notify(c, os.Interrupt)
	go func() {
		// Bloquear hasta que recibamos la señal SIGINT
		<-c
		common.Log.Infof(MSG_SERVER_SIGINT_ARRIVED)
		server.Stop()
	}()
}

// Permite servir a los clientes hasta que ocurre un error. Tiene una cantidad máxima de reintentos.
func (service *PeerService) Serve() {
	var err error = nil
	for range MAX_RETRY_SERVE {
		if err = service.Server.Serve(service.Listener); err != nil {
			common.Log.Debugf(MSG_FAILED_TO_SERVE, err)
		}
	}
	if err != nil {
		common.Log.Fatalf(MSG_FAILED_TO_SERVE, err)
	}
	common.Log.Infof(MSG_SERVER_STOPPED)
}
