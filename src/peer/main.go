package main

import (
	"fmt"
	"tp/peer/helpers"
)

const MESSAGE_START = "Starting node..."

func main() {
	helpers.Log.Info(MESSAGE_START)
	helpers.InitLogger()
	helpers.LoadConfig()
	fmt.Print(helpers.GetKey("gola"))
	fmt.Println("Hola Mundo")
}
