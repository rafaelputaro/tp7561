package main

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tp/common"
	filetransfer "tp/common/files_common/file_transfer"
	"tp/peer/helpers/file_manager/utils"
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
const MAX_RETRY_SERVE = 100
const MAX_INIT_SH_CTS_COUNT = 3

// Implementa la funcionalidad de grpc server para el par
type PeerService struct {
	Listener   net.Listener
	ServerGRPC *grpc.Server
	Receiver   *filetransfer.Receiver
}

// Retorna una nueva instancia de Peer Service lista para ser utilizada
func NewPeerService(peer *Peer) *PeerService {
	// GRPC service
	var lis net.Listener
	var err error = nil
	for range MAX_RETRY_LISTEN {
		lis, err = net.Listen("tcp", peer.Config.UrlGRPC)
		if err == nil {
			break
		}
		common.SleepBetweenRetries()
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
	// Escuchar llegada de archivos
	receiver, err := filetransfer.NewReceiver(
		peer.Config.UrlTCP,
		utils.GenerateIpfsUploadPath,
		receiveCallback(peer),
	)
	if err != nil {
		common.Log.Fatalf(MSG_FAILED_TO_LISTEN, err)
	}
	// Intercambio de contactos con boostrap
	shareContactsWithBootstrapNode(peer)
	// Detener el servidor cuando llega la señal SIGINT
	handleSigintSignal(server)
	// Servidor inicializado
	common.Log.Infof(MSG_SERVER_GRPC_STARTING)
	return &PeerService{
		Listener:   lis,
		ServerGRPC: server,
		Receiver:   receiver,
	}
}

// Intercambio frecuente de contactos con el bootstrapnode
func shareContactsWithBootstrapNode(peer *Peer) {
	go func() {
		count := 0
		for {
			common.SleepBetweenShareContactsShort()
			if count > MAX_INIT_SH_CTS_COUNT {
				common.SleepBetweenShareContactsLarge()
			} else {
				count++
			}
			if peer.NodeDHT.BucketTab.GetCountContacts() <= 3*peer.Config.NumberOfPairs/4 {
				peer.NodeDHT.ScheduleSndShCtsToBootstrapTask()
			}
		}
	}()
}

// Manejo de señal SIGINT
func handleSigintSignal(server *grpc.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		common.Log.Infof(MSG_SERVER_SIGINT_ARRIVED)
		server.Stop()
	}()
}

// Permite servir a los clientes hasta que ocurre un error. Tiene una cantidad máxima de reintentos.
func (service *PeerService) Serve() {
	var err error = nil
	for range MAX_RETRY_SERVE {
		if err = service.ServerGRPC.Serve(service.Listener); err != nil {
			common.Log.Debugf(MSG_FAILED_TO_SERVE, err)
		}
	}
	if err != nil {
		if !errors.Is(err, grpc.ErrServerStopped) {
			common.Log.Fatalf(MSG_FAILED_TO_SERVE, err)
		}
	}
	common.Log.Infof(MSG_SERVER_STOPPED)
}

// Callback ejecutado al recibir un archivo
func receiveCallback(peer *Peer) filetransfer.ReceiveCallback {
	return func(key []byte, fileName string) {
		// agregar archivo a la red de nodos
		peer.AddFileFromUploadDir(fileName)
	}
}
