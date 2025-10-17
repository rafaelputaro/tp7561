package filetransfer

import (
	"io"
	"net"
	"os"
	"tp/common"
	"tp/common/files_common/messages"
	"tp/common/keys"
)

// Es una entidad que administra la recepción de archivos
type Receiver struct {
	keyAndNameReceivedFiles map[string]string
	generatePath            func(string) string
}

// Retorna un nuevo recpetor listo para ser utilizado
func NewReceiver(selfUrl string, generatePath func(string) string) (*Receiver, error) {
	listener, err := net.Listen("tcp", selfUrl)
	if err != nil {
		common.Log.Debugf(MSG_ERROR_ON_STARTING_LISTENING, selfUrl, err)
		return nil, err
	}
	receiver := Receiver{
		keyAndNameReceivedFiles: map[string]string{},
		generatePath:            generatePath,
	}
	common.Log.Debugf(MSG_LISTENING_ON, selfUrl)
	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				common.Log.Debugf(MSG_ERROR_ACCEPTING_CONNECTION, err)
				continue
			}
			go receiver.handleConnection(conn)
		}
	}()
	return &receiver, nil
}

// Maneja una conexión guardando el archivo correspondiente
func (receiver *Receiver) handleConnection(conn net.Conn) {
	defer conn.Close()
	// Leer nombre del archivo
	buffer := make([]byte, BUFFER_SIZE)
	n, err := conn.Read(buffer)
	if err != nil {
		common.Log.Errorf(MSG_ERROR_READING_FILE_NAME, err)
		return
	}
	key, filename := parseKeyAndName(buffer[:n])
	// Crear archivo
	file, err := os.Create(receiver.generatePath(filename))
	if err != nil {
		common.Log.Errorf(messages.MSG_ERROR_CREATING_FILE, err)
		return
	}
	defer file.Close()
	common.Log.Debugf(MSG_RECEIVING_FILE, filename)
	// Recibir datos del archivo
	io.Copy(file, conn)
	common.Log.Debugf(MSG_FILE_RECEIVED_SUCCESSFULLY, filename)
	// Registro que el archivo ha sido recibido
	receiver.keyAndNameReceivedFiles[key] = filename
}

// obtiene la clave y el nombre del archivo <key><filename>
func parseKeyAndName(data []byte) (string, string) {
	dataS := string(data)
	return keys.KeyToHexString(keys.GetKey(dataS)), dataS
}
