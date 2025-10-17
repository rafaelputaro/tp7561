package communication

import (
	"net"
	"tp/common"
)

// Intenta conectarse como cliente a una url objetivo con retry
func ConnectAsClient(targetUrl string) (net.Conn, error) {
	var err error = nil
	var conn net.Conn = nil
	for range MAX_RETRIES_ON_CONNECT {
		conn, err = net.Dial("tcp", targetUrl)
		if err == nil {
			break
		}
		common.Log.Errorf(MSG_FAIL_ON_CONNECT_AS_CLIENT, err)
		common.SleepBetweenRetries()
	}
	return conn, err
}
