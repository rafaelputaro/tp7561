package helpers

import (
	"strconv"
	"tp/common/communication/url"
)

const CLIENT_HOST_PREFIX = "client-"
const BASE_PORT = 50051

// Retorna la url de un cliente
func GenerateURLClient(config Config, num int) string {
	host := CLIENT_HOST_PREFIX + strconv.Itoa(num)
	port := strconv.Itoa(BASE_PORT + config.NumberOfPairs + num)
	return url.GenerateURL(host, port)
}
