package helpers

import (
	"os"

	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("peer")

type TLoginFormatForKeys string

const (
	HEXA   TLoginFormatForKeys = "hex"
	BINARY TLoginFormatForKeys = "bin"
)

var LoginFormatForKeys = "bin"

// Inicializa el Log
func InitLogger() {
	format := logging.MustStringFormatter(
		`%{level:.5s} | %{shortfunc} | %{message}`,
	)
	LoginFormatForKeys = os.Getenv("LOGIN_FORMAT_FOR_KEYS")
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}
