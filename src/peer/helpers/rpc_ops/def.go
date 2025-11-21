package rpc_ops

const MSG_FAIL_ON_SEND_PING = "error sending ping: %v"
const MSG_FAIL_ON_SHARE_CONTACTS = "error on sharing contacts: %v"
const MSG_FAIL_ON_SEND_STORE = "error sending store message: %v"
const MSG_FAIL_ON_SEND_FIND_BLOCK = "error sending find block message: %v"
const MSG_PING_ATTEMPT = "ping attempt: %v | error: %v"
const MSG_SHARE_CONTACTS_ATTEMPT = "share contacts attempt: %v | error: %v"
const MSG_STORE_ATTEMPT = "store block attempt: %v | error: %v"
const MSG_FIND_BLOCK_ATTEMPT = "find block attempt: %v | error: %v"
const MAX_RETRIES_ON_PING = 10
const MAX_RETRIES_ON_SHARE_CONTACTS_RECIP = 10
const MAX_RETRIES_ON_STORE = 5
const MAX_RETRIES_ON_FIND_BLOCK = 20
