package rpc_ops

import (
	"tp/common"
	"tp/common/communication"
	"tp/common/contact"
	"tp/peer/helpers"
	"tp/protobuf/protoUtils"
)

// Ping con retry. En caso de no poder efectuar el ping retorna error
type PingOp func(config helpers.PeerConfig, contact contact.Contact) error

// Ping con retry. En caso de no poder efectuar el ping retorna error
func SndPing(config helpers.PeerConfig, contact contact.Contact) error {
	// conexión
	conn, client, ctx, cancel, err := communication.ConnectAsClientGRPC(contact.Url, communication.LogFatalOnFailConnectGRPC)
	if err == nil {
		defer conn.Close()
		defer cancel()
		// ping con retry
		for retry := range MAX_RETRIES_ON_PING {
			_, err = client.Ping(ctx, protoUtils.CreatePingOperands(config.Id, config.Url))
			if err != nil {
				common.Log.Infof(MSG_PING_ATTEMPT, retry, err)
				// esperar
				common.SleepBetweenRetries()
				continue
			}
			return nil
		}
		return err
	}
	common.Log.Errorf(MSG_FAIL_ON_SEND_PING, err)
	return err
}
