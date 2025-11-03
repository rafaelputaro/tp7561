package dht

import (
	"tp/common"
	"tp/common/contact"
	filetransfer "tp/common/files_common/file_transfer"
	"tp/common/keys"
	"tp/peer/helpers/file_manager"
	"tp/peer/helpers/file_manager/utils"
	peer_metrics "tp/peer/helpers/metrics"
)

const PREFIX_ADD_FILE = "up-"
const PREFIX_DOWNLOAD = "dow-"

// Genera el tag para una tarea de subida de archivo a la red de nodos
func generateAddFileTagFromFileName(fileName string) string {
	return generateAddFileTagFromKey(keys.GetKey(fileName))
}

// Genera el tag para una tarea de subida de archivo a la red de nodos
func generateAddFileTagFromKey(key []byte) string {
	return PREFIX_ADD_FILE + keys.KeyToHexString(key)
}

// Agrega la tarea de agregar un contacto a la bucket table. Se recomienda utilizarla
// para evitar posibles retrasos durante la actualización de la bucket table que impliquen
// enviar pings secundarios a otros contactos
func (node *Node) scheduleAddContactTask(contact contact.Contact) {
	node.TaskScheduler.AddTask(func() (string, bool) {
		node.BucketTab.AddContact(contact)
		return "", false
	})
}

// Agrega la tarea de agregar varios contactos a la bucket table. Se recomienda utilizarla
// para evitar posibles retrasos durante la actualización de la bucket table que impliquen
// enviar pings secundarios a otros contactos
func (node *Node) scheduleAddContactsTask(contacts []contact.Contact) {
	node.TaskScheduler.AddTask(func() (string, bool) {
		node.BucketTab.AddContacts(contacts)
		return "", false
	})
}

// Agrega la tarea de enviar ping a contactos para ser agregador a la bucket table
// en caso de encontrarse activos
func (node *Node) schedulePingAndAddContactsTask(contacts []contact.Contact) {
	for _, contact := range contacts {
		node.TaskScheduler.AddTask(func() (string, bool) {
			if node.SndPing(node.Config, contact) == nil {
				node.BucketTab.AddContact(contact)
			}
			return "", false
		})
	}
}

// Agrega la tarea de buscar el archivo
func (node *Node) scheduleGetFileTask(destUrl string, key []byte) error {
	return node.TaskScheduler.AddTask(func() (string, bool) {
		fileName, err := node.GetFileByKey(key)
		common.Log.Debugf(MSG_SENDING_FILE, fileName)
		if err == nil {
			filetransfer.SendFile(destUrl, fileName, utils.GenerateIpfsRestorePath(fileName))
			// Respaldar métrica
			peer_metrics.SetLastFileReturnedNumber(fileName)
			return "", false
		}
		common.Log.Debugf(MSG_ERROR_SEND_FILE, keys.KeyToLogFormatString(key), err)
		return "", false
	})
}

// Agrega la tarea se subir un archivo desde la carpeta upload
func (node *Node) scheduleAddFileFromUploadDirTask(fileName string) error {
	tag := generateAddFileTagFromFileName(fileName)
	return node.TaskScheduler.AddTaggedTask(func() (string, bool) {
		// si hay error, la tarea vuelve a intentar
		if err := file_manager.AddFileFromUploadDir(fileName, node.createSndBlockNeighbors()); err != nil {
			return tag, true
		}
		node.TaskScheduler.RemoveTaggedTask(tag)
		return tag, false
	}, tag)
}

// Agrega tarea de subir un archivo desde el espacio local
func (node *Node) scheduleAddFileFromInputDirTask(fileName string) error {
	tag := generateAddFileTagFromFileName(fileName)
	return node.TaskScheduler.AddTaggedTask(func() (string, bool) {
		if err := file_manager.AddFileFromInputDir(fileName, node.createSndBlockNeighbors()); err != nil {
			return tag, true
		}
		node.TaskScheduler.RemoveTaggedTask(tag)
		return tag, false
	}, tag)
}

// Agrega la tarea de envío de SndStore a un lote de contactos
func (node *Node) scheduleSndStoreTask(key []byte, fileName string, data []byte, contacts []contact.Contact) {
	for _, contact := range contacts {
		node.TaskScheduler.AddTask(func() (string, bool) {
			node.SndStore(node.Config, contact, key, fileName, data)
			return "", false
		})
	}
}
