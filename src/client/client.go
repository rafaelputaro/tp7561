package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
	"tp/client/helpers"
	client_metrics "tp/client/helpers/metrics"
	"tp/common"
	"tp/common/communication/url"
	"tp/common/files_common"
	filetransfer "tp/common/files_common/file_transfer"
	"tp/common/keys"
	rpc_ops_common "tp/common/rpc_ops"
	"tp/common/task_scheduler"
)

const MSG_ERROR_ON_CREATE_RECEIVER = "The file receiver could not be created"
const MSG_FILE_ADDED = "File added: %v | %v | %v"
const MSG_ERROR_FILE_NOT_UPLOADED = "the file has not yet been uploaded: %v"
const MSG_ERROR_ON_GET_FILE = "error on get file: %v"
const MSG_GET_FILE_ACCEPTED = "get file accepted: fileName: %v | key: %v"

type UploadData struct {
	Key         []byte
	UrlPeerDest string
}

type DownloadData struct {
	Key          []byte
	UrlPeerDest  string
	TimeReq      *time.Time
	TimeRcv      *time.Time
	DownloadTime int64 // tiempo que demoro la descarga
}

type Client struct {
	Config          helpers.Config
	UploadRegister  map[string]UploadData   // indexado por fileName
	DowloadRegister map[string]DownloadData // indexado por fileName
	Keys            map[string]string       // indexado por key como string y contiene el nombre del archivo
	Receiver        *filetransfer.Receiver
	TaskScheduler   task_scheduler.TaskScheduler
	Mutex           *sync.Mutex
}

// Retorna un cliente listo para ser utilizado
func NewClient() (*Client, error) {
	// configuración
	config := helpers.LoadConfig()
	// inicializar store
	helpers.InitStore(*config)
	// receptor de archivos por tcp
	receiver, err := filetransfer.NewReceiver(
		config.Url,
		func(fileName string) string {
			return helpers.GenerateDownloadPath(*config, fileName)
		},
		func([]byte, string) {},
	)
	if err != nil {
		common.Log.Errorf(MSG_ERROR_ON_CREATE_RECEIVER)
	}
	// iniciar servicio de métricas
	metrics := client_metrics.MetricsServiceInstance
	go func() {
		metrics.Serve()
	}()
	// retornar cliente
	return &Client{
		Config:          *config,
		UploadRegister:  map[string]UploadData{},
		DowloadRegister: map[string]DownloadData{},
		Keys:            map[string]string{},
		Receiver:        receiver,
		TaskScheduler:   *task_scheduler.NewTaskScheduler(),
		Mutex:           &sync.Mutex{},
	}, err
}

// Iniciar el cliente
func (client *Client) Start() {
	// esperar a que la mayoría de los pares se inicialicen intercambiando contactos
	common.SleepOnStart(10 * client.Config.NumberOfPairs)
	// agregar archivos
	client.addFiles()
	client.getFiles()
	for {
		common.SleepOnStart(10 * client.Config.NumberOfPairs)
	}
}

// Se detiene el cliente y sus servicios
func (client *Client) DisposeClient() {
	client.TaskScheduler.DisposeTaskScheduler()
}

// Intenta agregar todos los archivos del directorio de input a la red de nodos
func (client *Client) addFiles() error {
	return files_common.OpOverDir(helpers.GenerateInputFilePath(client.Config, ""),
		func(fileName string) error {
			keyS := keys.KeyToHexString(keys.GetKey(fileName))
			client.Keys[keyS] = fileName
			return client.scheduleAddFileTask(fileName)
		})
}

// Intenta enviar un archivo a la red de nodos
func (client *Client) addFile(fileName string) error {
	urlDest := selectPeer(client.Config.NumberOfPairs)
	key, err := rpc_ops_common.AddFile(urlDest, fileName, helpers.GenerateInputFilePath(client.Config, fileName))
	if err == nil {
		client.Mutex.Lock()
		defer client.Mutex.Unlock()
		common.Log.Debugf(MSG_FILE_ADDED, fileName, keys.KeyToHexString(key), urlDest)
		// agregar a los registros
		client.UploadRegister[fileName] = UploadData{
			Key:         key,
			UrlPeerDest: urlDest,
		}
		// agregar a las métricas
		client_metrics.MetricsServiceInstance.IncUploadedFileCount()
	}
	return err
}

// Intenta recuperar los archivos de la red de nodos
func (client *Client) getFiles() error {
	for _, fileName := range client.Keys {
		//
		client.scheduleGetFileTask(fileName)
		//client.getFile(fileName)
		common.SleepShort(client.Config.NumberOfPairs)
	}
	return nil
}

func (client *Client) getFile(fileName string) error {
	// Obtener datos de subida
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	// Si hay datos de subida intenta obtener el archivo
	if upReg, ok := client.UploadRegister[fileName]; ok {
		urlDest := selectPeerExcl(client.Config.NumberOfPairs, upReg.UrlPeerDest)
		_, _, errGet := rpc_ops_common.GetFile(client.Config.Url, urlDest, upReg.Key)
		if errGet != nil {
			msg := fmt.Sprintf(MSG_ERROR_ON_GET_FILE, errGet)
			common.Log.Debugf(msg)
			return errors.New(msg)
		}
		common.Log.Debugf(MSG_FILE_ADDED, fileName, keys.KeyToHexString(upReg.Key), upReg.UrlPeerDest)
		return nil
	}
	msg := fmt.Sprintf(MSG_ERROR_FILE_NOT_UPLOADED, fileName)
	common.Log.Debugf(msg)
	common.SleepShort(client.Config.NumberOfPairs)
	return errors.New(msg)
}

// Retorna la url de un peer aleatorio excluyendo a un par en particular
func selectPeerExcl(countPeers int, excl string) string {
	if countPeers > 1 {
		for {
			if generated := selectPeer(countPeers); generated != excl {
				return generated
			}
		}
	}
	return excl
}

// Retorna la url de un peer aleatorio
func selectPeer(countPeers int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	peerNumber := 2 + rand.Intn(countPeers-1)
	return url.GenerateURLPeer(peerNumber)
}
