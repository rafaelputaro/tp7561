package communication

import (
	"context"
	"time"

	"tp/common"
	"tp/protobuf/protopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//const MSG_SIGNIT_ARRIVED = "SIGNIT arrived. Stopping peer as gRPC Server"

// Se conecta a otro nodo como cliente grpc. En caso de conexión fallido luego de cierta cantidad
// de reintentos retorna un error
func ConnectAsClientGRPC(serverUrl string, callbackOnFailConn func(err error)) (*grpc.ClientConn, protopb.OperationsClient, context.Context, context.CancelFunc, error) {
	var err error = nil
	var conn *grpc.ClientConn = nil
	for range MAX_RETRIES_ON_CONNECT {
		conn, err = grpc.NewClient(serverUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			break
		}
		common.SleepBetweenRetries()
	}
	if err != nil {
		callbackOnFailConn(err)
		return conn, nil, nil, nil, err
	}
	c := protopb.NewOperationsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), MAX_TIMEOUT_ON_CONNECT*time.Second)
	return conn, c, ctx, cancel, err
}

func LogFatalOnFailConnectGRPC(err error) {
	common.Log.Fatalf(MSG_FAIL_ON_CONNECT_AS_CLIENT, err)
}

func LogErrorOnFailConnectGRPCError(err error) {
	common.Log.Errorf(MSG_FAIL_ON_CONNECT_AS_CLIENT, err)
}
