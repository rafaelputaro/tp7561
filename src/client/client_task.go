package main

import (
	"tp/common"
	"tp/common/keys"
)

const PREFIX_ADD_FILE = "add-file-"
const PREFIX_GET_FILE = "get-file-"
const MSG_RETRY_GET_FILE = "Retry to get file: %v"

// Genera el tag para una tarea de subida de archivo a la red de nodos
func generateAddFileTag(fileName string) string {
	return PREFIX_ADD_FILE + keys.KeyToHexString(keys.GetKey(fileName))
}

// Genera el tag para obtener un archivo desde la red de nodos
func generateGetFileTag(fileName string) string {
	return PREFIX_GET_FILE + fileName
}

// Agrega la tarea se subir un archivo desde la carpeta upload
func (client *Client) scheduleAddFileTask(fileName string) error {
	tag := generateAddFileTag(fileName)
	return client.TaskScheduler.AddTask(func() (string, bool) {
		// si hay error, la tarea se vuelve a intentar
		if client.addFile(fileName) != nil {
			common.SleepShort(client.Config.NumberOfPairs)
			return tag, true
		}
		common.SleepShort(client.Config.NumberOfPairs)
		return tag, false
	}, true, tag)
}

// Agrega la tarea de buscar el archivo
func (client *Client) scheduleGetFileTask(fileName string) error {
	tag := generateGetFileTag(fileName)
	return client.TaskScheduler.AddTask(func() (string, bool) {
		if err := client.getFile(fileName); err != nil {
			common.SleepShort(client.Config.NumberOfPairs)
			return tag, true
		}
		common.SleepShort(client.Config.NumberOfPairs)
		return tag, false
	}, false, tag)
}

// Volver a intentar obtener el archivo
func (client *Client) checkMustRetryGetFile(fileName string) {
	tag := generateGetFileTag(fileName)
	if !client.TaskScheduler.HasTag(tag) {
		common.Log.Debugf(MSG_RETRY_GET_FILE, fileName)
		client.scheduleGetFileTask(fileName)
	}
}

// Si el archivo ha sido agregado retorna true
func (client *Client) checkFileHasBeenAdded(fileName string) bool {
	tag := generateAddFileTag(fileName)
	return !client.TaskScheduler.HasTag(tag)
}
