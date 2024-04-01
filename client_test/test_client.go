package main

import (
	"fmt"
	"go-game-server/packet"
	"net"
	"time"
)

func main() {
	serverAddr := "127.0.0.1:42069" // Change this to the address and port of your server

	// Create three client goroutines
	for i := 1; i <= 3; i++ {
		go func(clientID int) {
			conn, err := net.Dial("tcp", serverAddr)
			if err != nil {
				fmt.Printf("Error connecting to server for client %d: %v\n", clientID, err)
				return
			}
			defer conn.Close()

			// Example: sending a login packet
			loginPacket := packet.Packet{
				ID:      packet.PLAYER_LOGIN_PACKET,
				Size:    5,                                       // Adjust the size according to your packet structure
				Payload: []byte(fmt.Sprintf("User%d", clientID)), // Example payload, change as needed
			}
			err = loginPacket.WritePacket(conn)
			if err != nil {
				fmt.Printf("Error sending login packet for client %d: %v\n", clientID, err)
				return
			}
			time.Sleep(time.Duration(clientID) * time.Second) // Sleep for a while before sending another packet

			// Example: sending player input packet
			//playerInputPacket := packet.Packet{
			//	ID:      packet.PLAYER_INPUT_PACKET,
			//	Size:    4,                  // Adjust the size according to your packet structure
			//	Payload: []byte{1, 0, 0, 1}, // Example payload representing player inputs, change as needed
			//}
			//err = playerInputPacket.WritePacket(conn)
			//if err != nil {
			//	fmt.Printf("Error sending player input packet for client %d: %v\n", clientID, err)
			//	return
			//}
			// You can continue sending packets here as needed
		}(i)
	}

	// Keep the main goroutine running
	select {}
}
