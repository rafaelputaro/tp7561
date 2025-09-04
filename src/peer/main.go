package main

import (
	"fmt"
	"tp/peer/helpers"
	//"tp/peer/helpers/rpc_ops/protobuf/protopb"
	//"google.golang.org/protobuf/proto"
)

const MESSAGE_START = "Starting node..."

func main() {
	helpers.Log.Info(MESSAGE_START)
	helpers.InitLogger()
	helpers.LoadConfig()
	fmt.Print(helpers.GetKey("gola"))
	fmt.Println("Hola Mundo")
	/*
		actor := protopb.Actor{
			Name:        proto.String("Dummys"),
			ProfilePath: proto.String("Dummy.jpg"),
			CountMovies: proto.Int64(0),
			ClientId:    proto.String("hola"),
			MessageId:   proto.Int64(0),
			SourceId:    proto.String("hola"),
			Eof:         proto.Bool(true),
		}
		fmt.Println("Actor: ", actor.GetClientId())*/
}
