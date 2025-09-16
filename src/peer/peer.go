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

func NewPeer(config helpers.PeerConfig) *Peer {
	peer := Peer{
		Config:  config,
		NodeDHT: *dht.NewNode(config, rpc_ops.SndPing, sndStore, rpc_ops.SndShareContactsRecip),
	}
	peer.GrpcService = *NewPeerService(&peer)
	return &peer
}

func (peer *Peer) Serve() {
	peer.GrpcService.Serve()
}

// Hace el procesamiento de la recepción de un ping desde el contacto parámetro e intenta agregarlo
// a la tabla de contactos.
func (peer *Peer) Ping(ctx context.Context, sourceContact *protopb.PingOperands) (*emptypb.Empty, error) {
	peer.NodeDHT.RcvPing(*contacts_queue.NewContact(sourceContact.GetSourceId(), sourceContact.GetSourceUrl()))
	return nil, nil
}

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

// Retorna los contactos de los nodos más cercanos a un targetId. Además hace el intento de
// agregar el contacto solicitante a la bucket_table

//func (node *Node) FindNode(contactSource contacts_queue.Contact, targetId []byte) []contacts_queue.Contact

// Si la target key se encuentra en el nodo retorna el valor de la misma, caso contrario retorna
// un error y la lista de los contactos más cercanos a la misma. Además hace el intento de
// agregar el contacto solicitante a la bucket_table

//func (node *Node) FindValue(contactSource contacts_queue.Contact, targetKey []byte) (string, []contacts_queue.Contact, error)

// Almacena la clave valor localmente y envía el menseja de store a los contactos más cercanos a la tabla.
// En caso de que la clave ya existía localmente retorna error. Por otro lado intenta agregar el contacto
// fuente en la tabla de contactos
//func (node *Node) Store(contactSource contacts_queue.Contact, key []byte, value string) error

func sndStore(config helpers.PeerConfig, contact contacts_queue.Contact, key []byte, value string) error {

	return nil
}
