package filetransfer

import (
	"os"
	"tp/common"
	"tp/common/communication"
	"tp/common/files_common/messages"
)

func SendFile(destUrl string, fileName string, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		common.Log.Errorf(messages.MSG_ERROR_OPENING_FILE, err)
	}
	defer file.Close()
	conn, err := communication.ConnectAsClient(destUrl)
	if err != nil {
		return err
	}
	defer conn.Close()
	// Enviar nombre del archivo
	_, err = conn.Write([]byte(filepath))
	if err != nil {
		common.Log.Errorf(messages.MSG_ERROR_SENDING_FILE_NAME, err)
	}
	/*
		// Enviar contenido del archivo
		_, err = io.Copy(conn, file)

			if err != nil {
				log.Fatal("Error al enviar archivo: ", err)
			}

		fmt.Println("Archivo enviado correctamente.")
	*/
	return nil
}

/*

// cliente.go
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

const BUFFER_SIZE = 1024

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: cliente <ruta-del-archivo>")
		return
	}

	filepath := os.Args
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Error al abrir archivo: ", err)
	}
	defer file.Close()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error al conectar al servidor: ", err)
	}
	defer conn.Close()

	// Enviar nombre del archivo
	_, err = conn.Write([]byte(filepath.Base(filepath)))
	if err != nil {
		log.Fatal("Error al enviar nombre del archivo: ", err)
	}

	// Enviar contenido del archivo
	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatal("Error al enviar archivo: ", err)
	}

	fmt.Println("Archivo enviado correctamente.")
}

*/
