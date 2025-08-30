package ring

import "strconv"

const BOOTSTRAP_NODE_ID = 1

const INVALID_ID = -1

const EMPTY_URL = ""

type Node struct {
	ID               int64
	Url              string
	NextID           int64
	NextUrl          string
	BootstrapNodeID  int64
	BootstrapNodeUrl int64
}

func getURL(ID int64) string {
	return "peer-" + strconv.FormatInt(ID, 10)
}

func (node *Node) Join() bool {
	// Soy el bootstrap
	if node.ID == node.BootstrapNodeID {
		return true
	}
	// Resolver en bootstrap node
	result, nextID, nextUrl := node.joinInBootstrapNode()
	if result {
		node.NextID = nextID
		node.NextUrl = nextUrl
		return true
	}
	return false
}

func (node *Node) joinInBootstrapNode() (bool, int64, string) {
	// @TODO comunicacion
	return true, INVALID_ID, EMPTY_URL
}

// Intenta resolver el join en el nodo actual.
// Retorna:
//  1. Si pudo unirlo al anillo retorna verdadero y el id del siguiente nodo
//  2. Si no pudo unirlo al anillo retorna falso y un id inválido
func (node *Node) JoinNewNode(newNodeID int64) (bool, int64, string) {
	if newNodeID <= node.ID {
		return false, INVALID_ID, EMPTY_URL
	}
	// Es el nodo inmediatamente siguiente al actual
	if newNodeID < node.NextID {
		node.NextID = newNodeID
		return true, node.NextID, node.NextUrl
	}
	// Resolver en el siguiente nodo
	return node.joinNewNodeInNext(newNodeID)
}

// Intenta resolver el join en el siguiente nodo.
// Retorna:
//  1. Si pudo unirlo al anillo retorna verdadero y el id del siguiente nodo
//  2. Si no pudo unirlo al anillo retorna falso y un id inválido
func (node *Node) joinNewNodeInNext(newNodeID int64) (bool, int64, string) {
	// @TODO comunicacion
	return true, INVALID_ID, EMPTY_URL
}
