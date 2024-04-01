package server

import (
	"fmt"
	"go-game-server/packet"
	"io"
	"net"
	"time"
)

const (
	NET_TCP        = "tcp"
	NET_TCP4       = "tcp4"
	NET_TCP6       = "tcp6"
	NET_UNIX       = "unix"
	NET_UNIXPACKET = "unixpacket"
)

type PlayerID string
type Player struct {
	ID         PlayerID
	Connection net.Conn
	Name       string
}

type GameServer struct {
	networkType    string
	address        string
	port           int
	fullAddress    string
	maxConnections int
	running        bool

	connectedPlayers map[PlayerID]*Player
}

func NewGameServer(t string, addr string, port int, maxConnections int) *GameServer {
	return &GameServer{
		networkType:    t,
		address:        addr,
		port:           port,
		fullAddress:    fmt.Sprintf("%s:%d", addr, port),
		maxConnections: maxConnections,
		running:        false,

		connectedPlayers: make(map[PlayerID]*Player),
	}
}

func (server *GameServer) StartServer() {
	listener, err := net.Listen(server.networkType, server.fullAddress)
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Server started, listening on %s...\n", server.fullAddress)

	// Start tick and fixedTick loops
	go server.tick()
	go server.fixedTick()

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		fmt.Printf("Client '%s' connecting...\n", conn.RemoteAddr())
		go server.handleConnection(conn)
	}
}

func (server *GameServer) tick() {
	for {
		// Print connected players
		fmt.Println("Connected players: ")
		for _, player := range server.connectedPlayers {
			if player == nil {
				continue
			}
			fmt.Printf("\t%v(%v)\n", player.Name, player.ID)
		}

		time.Sleep(time.Second) // Adjust the sleep duration as needed
	}
}

func (server *GameServer) fixedTick() {
	ticker := time.NewTicker(time.Second / 20)
	defer ticker.Stop()
	for {
		// Perform actions for fixedTick
		<-ticker.C
		//fmt.Println("Fixed Tick")
	}
}

func (server *GameServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Continuously read packets from the connection
	for {
		_packet, err := packet.ReadPacket(conn)
		if err != nil {
			switch err {
			// Handle client disconnect, EOF
			case io.EOF:
				server.handleClientDisconnect(conn)
				return
			default:
				fmt.Println("Error reading packet: ", err)
				return

			}
		}

		fmt.Printf("Client '%s' sent->\tID: %d,\t%d bytes,\t%v\n", conn.RemoteAddr(), _packet.ID, _packet.Size, _packet.Payload)

		switch _packet.ID {
		case packet.PLAYER_LOGIN_PACKET:
			server.handleLoginPacket(conn, _packet)
		case packet.WORLD_JOIN_REQUEST_PACKET:
			fmt.Println("WORLD_JOIN_REQUEST_PACKET-handler")
		case packet.PLAYER_INPUT_PACKET:
			server.handlePlayerInputPacket(conn, _packet)
		case packet.UPDATE_PLAYER_POSITION_PACKET:
			fmt.Println("UPDATE_PLAYER_POSITION_PACKET-handler")

		}
	}
}

func (server *GameServer) handleClientDisconnect(conn net.Conn) {
	server.connectedPlayers[PlayerID(conn.RemoteAddr().String())] = nil
	fmt.Printf("Client '%s' disconnected...\n", conn.RemoteAddr())
}

func (server *GameServer) handleLoginPacket(sender net.Conn, _packet packet.Packet) {
	newPlayer := Player{
		ID:         PlayerID(sender.RemoteAddr().String()),
		Connection: sender,
		Name:       string(_packet.Payload[:]),
	}

	server.connectedPlayers[newPlayer.ID] = &newPlayer
	fmt.Println("Player logged in: ", newPlayer)
}

func (server *GameServer) handlePlayerInputPacket(sender net.Conn, _packet packet.Packet) {
	inputUp := _packet.Payload[0]
	inputDown := _packet.Payload[1]
	inputLeft := _packet.Payload[2]
	inputRight := _packet.Payload[3]

	fmt.Printf(
		"%s - PlayerInput: [\n\tup:\t%v,\n\tdown:\t%v,\n\tleft:\t%v,\n\tright:\t%v,\n]\n",
		sender.RemoteAddr(), inputUp, inputDown, inputLeft, inputRight,
	)
}
