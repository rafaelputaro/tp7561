package main

import (
	"context"
	"tp/peer/dht"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
	"tp/peer/helpers/communication/rpc_ops"
	"tp/protobuf/protoUtils"
	"tp/protobuf/protopb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Peer struct {
	Config      helpers.PeerConfig
	NodeDHT     dht.Node
	GrpcService PeerService
	protopb.UnimplementedOperationsServer
}

// Retorna una instancia de peer lista para ser utilizada
func NewPeer(config helpers.PeerConfig) *Peer {
	peer := Peer{
		Config:  config,
		NodeDHT: *dht.NewNode(config, rpc_ops.SndPing, rpc_ops.SndStore, rpc_ops.SndShareContactsRecip),
	}
	peer.GrpcService = *NewPeerService(&peer)
	return &peer
}

// Inicia el servicio de atención de solicitudes rpc
func (peer *Peer) Serve() {
	peer.GrpcService.Serve()
}

// Hace el procesamiento de la recepción de un ping desde el contacto parámetro e intenta agregarlo
// a la tabla de contactos.
func (peer *Peer) Ping(ctx context.Context, sourceContact *protopb.PingOperands) (*emptypb.Empty, error) {
	peer.NodeDHT.RcvPing(*contacts_queue.NewContact(sourceContact.GetSourceId(), sourceContact.GetSourceUrl()))
	return nil, nil
}

// Intenta agregar los contactos recibidos a la bucket tabler, agregar el contacto fuente a la bucket table y
// retorna los contactos útiles para el contacto fuente
func (peer *Peer) ShareContactsReciprocally(ctx context.Context, sourceOperands *protopb.ShareContactsReciprocallyOperands) (*protopb.ShareContactsReciprocallyResults, error) {
	// parsear parámetros
	sourceContact, sourceContactList := protoUtils.ParseShareContactsReciprocallyOperands(sourceOperands)
	// agregar recomendados por la fuente y obtener recomendados
	selfContacts := peer.NodeDHT.RcvShareContactsReciprocally(sourceContact, sourceContactList)
	// agregar contactos que compartió la fuente
	return protoUtils.CreateShareContactsReciprocallyResults(selfContacts), nil
}

// Envía los contactos propios al bootstrap node esperando que el mismo retorne los contactos recomendados
// para la clave del presente nodo
func (peer *Peer) SndShareContactsToBootstrap() {
	peer.NodeDHT.SndShareContactsToBootstrap()
}

//func (node *Node) FindNode(contactSource contacts_queue.Contact, targetId []byte) []contacts_queue.Contact

// Si la target key se encuentra en el nodo retorna el valor de la misma, caso contrario retorna
// un error y la lista de los contactos más cercanos a la misma. Además hace el intento de
// agregar el contacto solicitante a la bucket_table

//func (node *Node) FindValue(contactSource contacts_queue.Contact, targetKey []byte) (string, []contacts_queue.Contact, error)

// Almacena la clave valor localmente y envía el menseja de store a los contactos más cercanos a la tabla.
// En caso de que la clave ya existía localmente retorna error. Por otro lado intenta agregar el contacto
// fuente en la tabla de contactos
func (peer *Peer) StoreBlock(ctx context.Context, operands *protopb.StoreBlockOperands) (*emptypb.Empty, error) {
	sourceContact, blockKey, blockName, data := protoUtils.ParseStoreBlockOperands(operands)
	peer.NodeDHT.RcvStore(*sourceContact, blockKey, blockName, data)
	return nil, nil
}

// Agrega un archivo local a la red de nodos del ipfs
func (peer *Peer) AddFile(fileName string) error {
	return peer.NodeDHT.AddFile(fileName)
}

func (peer *Peer) GetFile(fileName string) error {
	return peer.NodeDHT.GetFile(fileName)
}
