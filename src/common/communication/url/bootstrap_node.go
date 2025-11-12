package url

import (
	"os"
	"strconv"
	"tp/common"
	"tp/common/keys"
)

const BOOTSTRAP_NODE_DEFAULT = "peer-1"
const BOOTSTRAP_NODE_PORT = "50051"

var BootstrapNode = getBootrapNode()
var BootstrapNodeUrl = GenerateURL(BootstrapNode, BOOTSTRAP_NODE_PORT)
var BootstrapNodeID = keys.GetKey(BootstrapNode)
var IsBootstrapNodeSec = isBootstrapNodeSec()

func getBootrapNode() string {
	bootstrapNodeEnv := os.Getenv("BOOTSTRAP_NODE")
	if bootstrapNodeEnv == "" {
		return BOOTSTRAP_NODE_DEFAULT
	}
	return bootstrapNodeEnv
}

func isBootstrapNodeSec() bool {
	isBootstrapNodeSecS := os.Getenv("IS_BOOTSTRAP_NODE_SEC")
	parsed, err := strconv.ParseBool(isBootstrapNodeSecS)
	toReturn := false
	if err == nil {
		toReturn = parsed
	}
	if parsed {
		common.Log.Debugf("is a bootstrap node sec ")
	} else {
		common.Log.Debugf("is not a bootstrap node sec")
	}
	return toReturn
}
