package filetransfer

import (
	"fmt"
	"io"
	"os"
	"tp/common"
	"tp/common/communication"
	"tp/common/files_common/messages"
)

// Enviar archivo a una url dada
func SendFile(destUrl string, fileName string, filePath string) error {
	// Abrir archivo
	file, err := os.Open(filePath)
	if err != nil {
		common.Log.Errorf(messages.MSG_ERROR_OPENING_FILE, err)
		return err
	}
	defer file.Close()
	// Conectarse como cliente
	conn, err := communication.ConnectAsClient(destUrl)
	if err != nil {
		return err
	}
	defer conn.Close()
	// Enviar nombre del archivo
	fileNameToSend := fmt.Sprintf("%v\n", fileName)
	_, err = conn.Write([]byte(fileNameToSend))
	if err != nil {
		common.Log.Errorf(messages.MSG_ERROR_SENDING_FILE_NAME, err)
	}
	// Enviar contenido del archivo
	_, err = io.Copy(conn, file)
	if err != nil {
		common.Log.Errorf(MSG_ERROR_SENDING_FILE, err)
		return err
	}
	common.Log.Debugf(MSG_FILE_SENT_SUCCESSFULLY, fileName)
	return nil
}
