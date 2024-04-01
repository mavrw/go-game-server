package main

import (
	"fmt"
	"go-game-server/server"
)

var IP_ADDRESS = "127.0.0.1"
var PORT = 42069
var MAX_CONNECTIONS = 128

func main() {
	gameServer := server.NewGameServer(server.NET_TCP4, IP_ADDRESS, PORT, MAX_CONNECTIONS)

	fmt.Println("Starting server...")
	gameServer.StartServer()
	fmt.Println("Exiting...")

}
