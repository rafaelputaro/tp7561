package dht

import (
	"bytes"
	"tp/peer/dht/bucket_table"
	"tp/peer/helpers"
	"tp/peer/helpers/rpc_ops"
)

const BOOTSTRAP_NODE_NAME = "peer-1"
const BOOTSTRAP_NODE_PORT = "50051"
const BOOTSTRAP_NODE_HOST = "peer-1"

var BootstrapNodeUrl = helpers.GenerateURL(BOOTSTRAP_NODE_HOST, BOOTSTRAP_NODE_PORT)
var BootstrapNodeID = helpers.GetKey(BOOTSTRAP_NODE_NAME)

type Node struct {
	ID          []byte
	Url         string
	BucketTab   bucket_table.BucketTable
	KeyValueTab KeyValueTable
	ping        rpc_ops.PingOp
	// cache
}

// Retorna una nueva instancia de nodo
func NewNode(config helpers.PeerConfig, ping rpc_ops.PingOp) *Node {
	node := &Node{
		ID:          helpers.GetKey(config.Name),
		Url:         config.Url,
		KeyValueTab: *NewKeyValueTable(),
		ping:        ping,
	}
	return node
}

// Retorna verdadero si la instancia el bootstrap node
func (node *Node) IsBootstrapNode() bool {
	return bytes.Equal(node.ID, BootstrapNodeID)
}

//Ping

//FindNode

//FindValue

//Store
