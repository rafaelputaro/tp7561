package url

import (
	"strconv"
)

const PEER_HOST_PREFIX = "peer-"
const BASE_PORT = 50051

// Retorna la url de un peer
func GenerateURLPeer(num int) string {
	host := PEER_HOST_PREFIX + strconv.Itoa(num)
	port := strconv.Itoa(BASE_PORT)
	return GenerateURL(host, port)
}

// Retorna la Url en base al host y el puerto
func GenerateURL(host string, port string) string {
	return host + ":" + port
}
