package connection

import (
	"context"
	"time"

	"tp/peer/common/helpers"
	"tp/peer/common/protobuf/protopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const MAX_RETRIES_ON_CONNECT = 20
const MAX_TIMEOUT_ON_CONNECT = 60
const MSG_FAIL_ON_CONNECT_AS_CLIENT = "error trying to connect as a client: %v"
const MSG_SIGNIT_ARRIVED = "SIGNIT arrived. Stopping peer as gRPC Server"

// Se conecta a otro nodo como cliente grpc. En caso de conexión fallido luego de cierta cantidad
// de reintentos retorna un error
func ConnectAsClient(serverUrl string, callbackOnFailConn func(err error)) (*grpc.ClientConn, protopb.OperationsClient, context.Context, context.CancelFunc, error) {
	var err error = nil
	var conn *grpc.ClientConn = nil
	for range MAX_RETRIES_ON_CONNECT {
		conn, err = grpc.NewClient(serverUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			break
		}
		helpers.SleepBetweenRetries()
	}
	if err != nil {
		callbackOnFailConn(err)
		return conn, nil, nil, nil, err
	}
	c := protopb.NewOperationsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), MAX_TIMEOUT_ON_CONNECT*time.Second)
	return conn, c, ctx, cancel, err
}

/*
func ConnectAsServer(peer common.Peer) {

	lis, err := net.Listen("tcp", peer.Config.Url) //":"+config.Port)
	if err != nil {
		helpers.Log.Fatalf("failed to listen: %v", err)
	}
	// New server
	server := grpc.NewServer()
	protopb.RegisterOperationsServer(server, &peer)
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
		// Se bloque hasta recibir una señal SIGINT
		<-c
		helpers.Log.Infof(MSG_SIGNIT_ARRIVED)
		server.Stop()
	}()
}

*/

func LogFatalOnFailConnect(err error) {
	helpers.Log.Fatalf(MSG_FAIL_ON_CONNECT_AS_CLIENT, err)
}

func LogErrorOnFailConnectError(err error) {
	helpers.Log.Errorf(MSG_FAIL_ON_CONNECT_AS_CLIENT, err)
}
