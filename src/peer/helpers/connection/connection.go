package connection

import (
	"context"
	"time"

	"tp/peer/helpers"
	"tp/peer/protobuf/protopb"

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

func LogFatalOnFailConnect(err error) {
	helpers.Log.Fatalf(MSG_FAIL_ON_CONNECT_AS_CLIENT, err)
}

func LogErrorOnFailConnectError(err error) {
	helpers.Log.Errorf(MSG_FAIL_ON_CONNECT_AS_CLIENT, err)
}
