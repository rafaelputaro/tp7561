package main

import (
	"context"
	"fmt"
	"tp/common"
	"tp/common/contact"
	"tp/common/files_common/messages"
	"tp/common/keys"
	"tp/peer/dht"
	"tp/peer/helpers"
	"tp/peer/helpers/file_manager"
	proto_utils_peer "tp/peer/helpers/proto_utils"
	"tp/peer/helpers/rpc_ops"
	"tp/protobuf/protoUtils"
	"tp/protobuf/protopb"

	"google.golang.org/protobuf/types/known/emptypb"
)

// Contiene la funcionalidad de IPFS y DHT
type Peer struct {
	Config      helpers.PeerConfig
	NodeDHT     dht.Node
	GrpcService PeerService
	protopb.UnimplementedOperationsServer
}

// Retorna una instancia de peer lista para ser utilizada
func NewPeer(config helpers.PeerConfig) *Peer {
	peer := Peer{
		Config: config,
		NodeDHT: *dht.NewNode(
			config, rpc_ops.SndPing,
			rpc_ops.SndStore,
			rpc_ops.SndShareContactsRecip,
			rpc_ops.SndFindBlock),
	}
	peer.GrpcService = *NewPeerService(&peer)
	return &peer
}

// Liberar recursos
func (peer *Peer) DisposePeer() {
	peer.NodeDHT.DisposeNode()
}

// Inicia el servicio de atención de solicitudes rpc
func (peer *Peer) Serve() {
	peer.GrpcService.Serve()
}

// Agrega un archivo local a la red de nodos del ipfs
func (peer *Peer) DoAddFile(fileName string) error {
	return peer.NodeDHT.AddFile(fileName)
}

// Agrega una parte de archivo al peer por rpc. En caso de ser la última parte inicia el proceso
// de subida a la red
func (peer *Peer) AddFile(ctx context.Context, fileOpers *protopb.AddFileOpers) (*protopb.AddFileRes, error) {
	fileName, part, data, endFile := protoUtils.ParseAddFileOperands(fileOpers)
	var key []byte = keys.GetNullKey()
	var err error = nil
	// Si no existe la key en la tabla se procede a guardar las partes
	if !peer.existKeyOrFileExistInUploadDir(fileName) {
		restored, errSt := file_manager.StoreUploadFilePart(fileName, part, data, endFile)
		if restored {
			// programar DoAddFile y retornar clave
			key = keys.GetKey(fileName)
			//peer.DoAddFile(fileName) Modificarlo para que tome desde la carpeta upload
		}
		err = errSt
	} else {
		common.Log.Debugf(fmt.Sprintf("%v: %v", messages.MSG_ERROR_FILE_EXIST, fileName))
		key = keys.GetKey(fileName)
	}
	return protoUtils.CreateAddFileResults(key), err
}

func (peer *Peer) GetFile(fileName string) error {
	return peer.NodeDHT.GetFile(fileName)
}

func (peer *Peer) existKeyOrFileExistInUploadDir(fileName string) bool {
	// Busca en la bucket table
	if peer.NodeDHT.ExistsFileLocally(fileName) {
		return true
	}
	// Busca en el directorio upload
	return file_manager.FileExistInUpload(fileName)
}

// Hace el procesamiento de la recepción de un ping desde el contacto parámetro e intenta agregarlo
// a la tabla de contactos.
func (peer *Peer) Ping(ctx context.Context, sourceContact *protopb.PingOperands) (*emptypb.Empty, error) {
	peer.NodeDHT.RcvPing(*contact.NewContact(sourceContact.GetSourceId(), sourceContact.GetSourceUrl()))
	return nil, nil
}

// Intenta agregar los contactos recibidos a la bucket tabler, agregar el contacto fuente a la bucket table y
// retorna los contactos útiles para el contacto fuente
func (peer *Peer) ShCtsReciprocally(ctx context.Context, sourceOperands *protopb.ShCtsRecipOpers) (*protopb.ShCtsRecipRes, error) {
	// parsear parámetros
	sourceContact, sourceContactList := protoUtils.ParseShareContactsReciprocallyOperands(sourceOperands)
	// agregar recomendados por la fuente y obtener recomendados
	selfContacts := peer.NodeDHT.RcvShCtsRecip(sourceContact, sourceContactList)
	// agregar contactos que compartió la fuente
	return protoUtils.CreateShareContactsReciprocallyResults(selfContacts), nil
}

// Envía los contactos propios al bootstrap node esperando que el mismo retorne los contactos recomendados
// para la clave del presente nodo
func (peer *Peer) SndShCtsToBootstrap() {
	peer.NodeDHT.SndShCtsToBootstrap()
}

// Almacena la clave valor localmente y envía el menseja de store a los contactos más cercanos a la tabla.
// En caso de que la clave ya existía localmente retorna error. Por otro lado intenta agregar el contacto
// fuente en la tabla de contactos
func (peer *Peer) StoreBlock(ctx context.Context, operands *protopb.StoreBlockOpers) (*emptypb.Empty, error) {
	sourceContact, blockKey, blockName, data := protoUtils.ParseStoreBlockOperands(operands)
	peer.NodeDHT.RcvStore(*sourceContact, blockKey, blockName, data)
	return nil, nil
}

// Si la target key se encuentra en el nodo retorna el valor de la misma, caso contrario retorna
// un error y la lista de los contactos más cercanos a la misma. Además hace el intento de
// agregar el contacto solicitante a la bucket_table
func (peer *Peer) FindBlock(ctx context.Context, operands *protopb.FindBlockOpers) (*protopb.FindBlockRes, error) {
	sourceContact, key := proto_utils_peer.ParseFindBlockOperands(operands)
	fileName, data, contacts, err := peer.NodeDHT.RcvFindBlock(sourceContact, key)
	return proto_utils_peer.CreateFindBlockResults(fileName, data, contacts), err
}
