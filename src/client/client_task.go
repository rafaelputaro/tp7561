package main

import (
	"tp/common"
)

const PREFIX_ADD_FILE = "add-file-"
const PREFIX_GET_FILE = "get-file-"
const MSG_RETRY_GET_FILE = "Retry to get file: %v"

// Genera el tag para obtener un archivo desde la red de nodos
func generateGetFileTag(fileName string) string {
	return PREFIX_GET_FILE + fileName
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
