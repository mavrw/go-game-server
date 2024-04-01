package packet

import (
	"encoding/binary"
	"net"
)

type PacketID uint32

// Packet ID types
const (
	PLAYER_LOGIN_PACKET PacketID = iota + 1
	WORLD_JOIN_REQUEST_PACKET
	PLAYER_INPUT_PACKET
	UPDATE_PLAYER_POSITION_PACKET
)

type Packet struct {
	ID      PacketID
	Size    uint32
	Payload []byte
}

func ReadPacket(conn net.Conn) (Packet, error) {
	// Read packet ID from first 4 bytes
	idBuffer := make([]byte, 4)
	_, err := conn.Read(idBuffer)
	if err != nil {
		return Packet{}, err
	}
	packetID := binary.BigEndian.Uint32(idBuffer)

	// Read payload size from next 4 bytes
	sizeBuffer := make([]byte, 4)
	_, err = conn.Read(sizeBuffer)
	if err != nil {
		return Packet{}, err
	}
	payloadSize := binary.BigEndian.Uint32(sizeBuffer)

	// Read payload
	payloadBuffer := make([]byte, payloadSize)
	_, err = conn.Read(payloadBuffer)
	if err != nil {
		return Packet{}, err
	}

	return Packet{
		ID:      PacketID(packetID),
		Size:    payloadSize,
		Payload: payloadBuffer,
	}, nil
}

func (packet *Packet) WritePacket(conn net.Conn) error {
	// Write packet ID
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, uint32(packet.ID))
	_, err := conn.Write(idBytes)
	if err != nil {
		return err
	}

	// Write packet size
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, packet.Size)
	_, err = conn.Write(sizeBytes)
	if err != nil {
		return err
	}

	// Write payload
	_, err = conn.Write(packet.Payload)
	if err != nil {
		return err
	}

	return nil
}
