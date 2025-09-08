package main

import (
	"context"
	"log"
	"time"
	"tp/peer/dht"
	"tp/peer/dht/bucket_table/contacts_queue"
	"tp/peer/helpers"
	"tp/peer/protobuf/protopb"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Peer struct {
	Config  helpers.PeerConfig
	NodeDHT dht.Node
	protopb.UnimplementedOperationsServer
}

func NewPeer(config helpers.PeerConfig) *Peer {
	peer := Peer{
		Config:  config,
		NodeDHT: *dht.NewNode(config, sndPing, sndStore),
	}
	return &peer
}

// @Todo modificar ping para que envíe los datos del contacto mediante rpc
func (peer *Peer) Ping(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	//peer.NodeDHT.Ping()
	helpers.Log.Debugf("Pong desde: %v", helpers.KeyToLogFormatString(peer.Config.Id))
	return nil, nil
}

func (peer *Peer) PingToBootstrap() {
	peer.NodeDHT.PingToBootstrap()
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

// @TODO agrega retry, mejorar todo el método, etc
func sndPing(config helpers.PeerConfig, contact contacts_queue.Contact) error {
	// Set up a connection to the gRPC server @TODO ARREGLAR ESTO PARA QUE NO ESTE DEPRECADO
	conn, err := grpc.Dial(contact.Url, grpc.WithInsecure())
	if err != nil {
		helpers.Log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create gRPC stub
	c := protopb.NewOperationsClient(conn)

	// Golang context pattern used to handle timeouts against the server.
	// Defined with a 5 seconds timeout but not used in the example
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.Ping(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not call MyMethod: %v", err)
	}

	return nil
}
