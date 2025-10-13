package helpers

import "tp/common/keys"

const BOOTSTRAP_NODE_NAME = "peer-1"
const BOOTSTRAP_NODE_PORT = "50051"
const BOOTSTRAP_NODE_HOST = "peer-1"

var BootstrapNodeUrl = GenerateURL(BOOTSTRAP_NODE_HOST, BOOTSTRAP_NODE_PORT)
var BootstrapNodeID = keys.GetKey(BOOTSTRAP_NODE_NAME)
