package rpc_ops

type PingOp func(url string) bool

type FindNodeOp func(id []byte)

type FindValueOp func()

type StoreOp func()
