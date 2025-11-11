package url

import (
	"os"
	"tp/common/keys"
)

const BOOTSTRAP_NODE_DEFAULT = "peer-1"
const BOOTSTRAP_NODE_PORT = "50051"

var bootstrapNode = getBootrapNode()
var BootstrapNodeUrl = GenerateURL(bootstrapNode, BOOTSTRAP_NODE_PORT)
var BootstrapNodeID = keys.GetKey(bootstrapNode)

func getBootrapNode() string {
	bootstrapNodeEnv := os.Getenv("BOOTSTRAP_NODE")
	if bootstrapNodeEnv == "" {
		return BOOTSTRAP_NODE_DEFAULT
	} else {
		return bootstrapNodeEnv
	}
}
