package filetransfer

import (
	"io"
	"net"
	"os"
	"strings"
	"tp/common"
	"tp/common/files_common/messages"
	"tp/common/keys"
)

// Función a ejectura luego de guardar los datos
type ReceiveCallback func(key []byte, fileName string)

// Es una entidad que administra la recepción de archivos
type Receiver struct {
	keyAndNameReceivedFiles map[string]string
	generatePath            func(string) string
	callback                ReceiveCallback
}

// Retorna un nuevo recpetor listo para ser utilizado.
// callback(<key>, <fileName>)
func NewReceiver(selfUrl string, generatePath func(string) string, callback ReceiveCallback) (*Receiver, error) {
	listener, err := net.Listen("tcp", selfUrl)
	if err != nil {
		common.Log.Debugf(MSG_ERROR_ON_STARTING_LISTENING, selfUrl, err)
		return nil, err
	}
	receiver := Receiver{
		keyAndNameReceivedFiles: map[string]string{},
		generatePath:            generatePath,
		callback:                callback,
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
	key, keyS, filename, data := parseKeyNameAndData(buffer[:n])
	// Crear archivo
	file, err := os.Create(receiver.generatePath(filename))
	if err != nil {
		common.Log.Errorf(messages.MSG_ERROR_CREATING_FILE, err)
		return
	}
	defer file.Close()
	common.Log.Debugf(MSG_RECEIVING_FILE, filename)
	// Escribir los bytes sobrantes de lectura nombre
	file.Write(data)
	// Recibir datos del archivo
	io.Copy(file, conn)
	common.Log.Debugf(MSG_FILE_RECEIVED_SUCCESSFULLY, filename)
	// Registro que el archivo ha sido recibido
	receiver.keyAndNameReceivedFiles[keyS] = filename
	receiver.callback(key, filename)
}

// obtiene la clave y el nombre del archivo <key><filename><datos restantes>
func parseKeyNameAndData(data []byte) ([]byte, string, string, []byte) {
	dataS := string(data)
	fileName, _, _ := strings.Cut(dataS, "\n")
	key := keys.GetKey(fileName)
	keyS := keys.KeyToHexString(key)
	return key, keyS, fileName, data[len(fileName)+1:]
}
