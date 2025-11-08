package dht

import (
	"fmt"
	"strconv"
	"time"
	"tp/common"
	"tp/common/communication/url"
	"tp/common/contact"
	filetransfer "tp/common/files_common/file_transfer"
	"tp/common/keys"
	"tp/peer/helpers/file_manager"
	"tp/peer/helpers/file_manager/utils"
	peer_metrics "tp/peer/helpers/metrics"
)

const PREFIX_ADD_CONTACT = "add-contact-"
const PREFIX_PING_AND_ADD_CONTACT = "ping-add-contact-"
const PREFIX_ADD_CONTACTS = "add-contacts-"
const PREFIX_ADD_FILE_FROM_INPUT = "up-from-inputDir-"
const PREFIX_ADD_FILE_FROM_UPLOAD = "up-from-uploadDir-"
const PREFIX_GET_FILE = "get-file-"
const PREFIX_SEND_FILE = "send-file-"
const PREFIX_SND_STORE = "snd-store-"
const PREFIX_SND_SH_CTS_BOOTSTRAP_NODE = "snd-sh-cts-bn-"

// Retorna un tag basado en el tiempo y el prefijo
func generateTimeNanoTag(prefix string) string {
	return fmt.Sprintf("%v%v", prefix, strconv.FormatInt(time.Now().UnixNano(), 10))
}

// Generar el tag para una tarea de agregar contacto
func generateAddContactTag() string {
	return generateTimeNanoTag(PREFIX_ADD_CONTACT)
}

// Genera el tag para una tarea de hacer un ping y luego agregar el contacto
func generatePingAndAddContactTag() string {
	return generateTimeNanoTag(PREFIX_PING_AND_ADD_CONTACT)
}

// Genera el tag para una tarea de agregar contactos.
func generateAddContactsTag() string {
	return generateTimeNanoTag(PREFIX_ADD_CONTACTS)
}

// Genera el tag para una tarea de subida de archivo a la red de nodos desde upload
func generateAddFileFromUploadTag(fileName string) string {
	return PREFIX_ADD_FILE_FROM_UPLOAD + keys.KeyToHexString(keys.GetKey(fileName))
}

// Genera el tag para una tarea de subida de archivo a la red de nodos desde input
func generateAddFileFromInputTag(fileName string) string {
	return PREFIX_ADD_FILE_FROM_INPUT + keys.KeyToHexString(keys.GetKey(fileName))
}

// Genera el tag para una tarea de subida de archivo a la red de nodos
func generateSndStoreFromInputTag(fileName string) string {
	return generateTimeNanoTag(PREFIX_SND_STORE + keys.KeyToHexString(keys.GetKey(fileName)))
}

// Genera el tag para obtener un archivo desde la red de nodos
func generateGetFileTag(destUrl string, key []byte) string {
	return PREFIX_GET_FILE + destUrl + "-" + keys.KeyToHexString(key)
}

// Genera el tag para el envío de un archivo a un cliente
func generateSendFileTag(destUrl string, key []byte) string {
	return PREFIX_SEND_FILE + destUrl + "-" + keys.KeyToHexString(key)
}

// Genera el tag para una tarea de compartir contactos con bootstrap node
func generateSndShCtsToBootstrapTag() string {
	return generateTimeNanoTag(PREFIX_SND_SH_CTS_BOOTSTRAP_NODE)
}

// Agrega la tarea de agregar un contacto a la bucket table. Se recomienda utilizarla
// para evitar posibles retrasos durante la actualización de la bucket table que impliquen
// enviar pings secundarios a otros contactos
func (node *Node) scheduleAddContactTask(contact contact.Contact) {
	tag := generateAddContactTag()
	node.TaskScheduler.AddTask(func() (string, bool) {
		// si hay error, la tarea se vuelve a intentar
		if node.BucketTab.AddContact(contact) != nil {
			return tag, true
		}
		return "", false
	}, tag)
}

// Agrega la tarea de agregar varios contactos a la bucket table. Se recomienda utilizarla
// para evitar posibles retrasos durante la actualización de la bucket table que impliquen
// enviar pings secundarios a otros contactos
func (node *Node) scheduleAddContactsTask(contacts []contact.Contact) {
	tag := generateAddContactsTag()
	node.TaskScheduler.AddTask(func() (string, bool) {
		// si hay error, la tarea se vuelve a intentar
		if node.BucketTab.AddContacts(contacts) != nil {
			return tag, true
		}
		return tag, false
	}, tag)
}

// Agrega la tarea de enviar ping a contactos para ser agregador a la bucket table
// en caso de encontrarse activos
func (node *Node) schedulePingAndAddContactsTask(contacts []contact.Contact) {
	for _, contact := range contacts {
		tag := generatePingAndAddContactTag()
		node.TaskScheduler.AddTask(func() (string, bool) {
			// si hay error, la tarea se vuelve a intentar
			if node.SndPing(node.Config, contact) != nil {
				return tag, true
			}
			node.scheduleAddContactTask(contact)
			return tag, false
		}, tag)
	}
}

// Agrega la tarea de buscar el archivo
func (node *Node) scheduleGetFileTask(destUrl string, key []byte) error {
	tag := generateGetFileTag(destUrl, key)
	return node.TaskScheduler.AddTask(func() (string, bool) {
		fileName, err := node.GetFileByKey(key)
		common.Log.Debugf(MSG_SENDING_FILE, fileName)
		if err == nil {
			// Programar retorno de archivo
			node.scheduleSendFile(destUrl, key, fileName)
			return tag, false
		}
		common.Log.Errorf(MSG_ERROR_GET_FILE, keys.KeyToLogFormatString(key), err)
		return tag, true
	}, tag)
}

// Agrega la tarea de envío de un archivo
func (node *Node) scheduleSendFile(destUrl string, key []byte, fileName string) error {
	tag := generateSendFileTag(destUrl, key)
	return node.TaskScheduler.AddTask(func() (string, bool) {
		common.Log.Debugf(MSG_SENDING_FILE, fileName)
		if err := filetransfer.SendFile(destUrl, fileName, utils.GenerateIpfsRestorePath(fileName)); err != nil {
			common.Log.Errorf(MSG_ERROR_SEND_FILE, keys.KeyToLogFormatString(key), err)
			return tag, true
		}
		// Respaldar métrica
		peer_metrics.SetLastFileReturnedNumber(fileName)
		return tag, false
	}, tag)
}

// Retorna verdadero si se encuentra pendiente la búsqueda o el envío de un archivo
func (node *Node) checkSendFilePending(destUrl string, key []byte) bool {
	tagGet := generateGetFileTag(destUrl, key)
	if node.TaskScheduler.HasTag(tagGet) {
		return true
	}
	tagSnd := generateSendFileTag(destUrl, key)
	return node.TaskScheduler.HasTag(tagSnd)
}

// Agrega la tarea se subir un archivo desde la carpeta upload
func (node *Node) scheduleAddFileFromUploadDirTask(fileName string) error {
	tag := generateAddFileFromUploadTag(fileName)
	return node.TaskScheduler.AddTask(func() (string, bool) {
		// si hay error, la tarea se vuelve a intentar
		if file_manager.AddFileFromUploadDir(fileName, node.createSndBlockNeighbors()) != nil {
			return tag, true
		}
		return tag, false
	}, tag)
}

// Agrega tarea de subir un archivo desde el espacio local
func (node *Node) scheduleAddFileFromInputDirTask(fileName string) error {
	tag := generateAddFileFromInputTag(fileName)
	return node.TaskScheduler.AddTask(func() (string, bool) {
		// si hay error, la tarea se vuelve a intentar
		if file_manager.AddFileFromInputDir(fileName, node.createSndBlockNeighbors()) != nil {
			return tag, true
		}
		return tag, false
	}, tag)
}

// Agrega la tarea de envío de SndStore a un lote de contactos
func (node *Node) scheduleSndStoreTask(key []byte, fileName string, data []byte, contacts []contact.Contact) {
	for _, contact := range contacts {
		tag := generateSndStoreFromInputTag(fileName)
		node.TaskScheduler.AddTask(func() (string, bool) {
			// si hay error, la tarea se vuelve a intentar
			if node.SndStore(node.Config, contact, key, fileName, data) != nil {
				return tag, true
			}
			return tag, false
		}, tag)
	}
}

// Agrega la tarea de compartir contactos con el boostrap node
func (node *Node) ScheduleSndShCtsToBootstrapTask() {
	tag := generateSndShCtsToBootstrapTag()
	node.TaskScheduler.AddTask(func() (string, bool) {
		if !node.IsBootstrapNode() {
			contactBoostrapNode := contact.NewContact(url.BootstrapNodeID, url.BootstrapNodeUrl)
			// agregar bootstrap node a contactos
			if node.SndShCts(*contactBoostrapNode) == nil {
				node.scheduleAddContactTask(*contactBoostrapNode)
				return tag, false
			}
			return tag, true
		}
		return tag, false
	}, tag)
}
