package sphere

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	// PacketEventMessage denotes a regular text data message.
	PacketEventMessage = 0
	// PacketEventChannel denotes a channel request.
	PacketEventChannel = 1
	// PacketEventSubscribe denotes a subscribe request.
	PacketEventSubscribe = 6
	// PacketEventUnsubscribe denotes an unsubscribe request.
	PacketEventUnsubscribe = 7
	// PacketEventPing denotes an ping message.
	PacketEventPing = 8
	// PacketEventPong denotes an pong message.
	PacketEventPong = 9
)

// NewPacket creates new packet instance
func NewPacket(event int, data string) *Packet {
	return &Packet{event, data, nil}
}

// ParsePacket parses string and return packet
func ParsePacket(str string) (*Packet, error) {
	if str == "" {
		return nil, errors.New("packet is empty")
	}
	code := string(str[0])
	i, err := strconv.Atoi(code)
	if err != nil {
		return nil, errors.New("packet is invalid")
	}
	if i != PacketEventMessage && i != PacketEventSubscribe && i != PacketEventUnsubscribe && i != PacketEventPing && i != PacketEventPong {
		return nil, errors.New("packet is invalid")
	}
	data := str[1:len(str)]
	if i == PacketEventSubscribe || i == PacketEventUnsubscribe || i == PacketEventMessage {
		if data == "" {
			return nil, errors.New("packet is invalid")
		}
	}
	return NewPacket(i, data), nil
}

// Packet indicates the data of the message
type Packet struct {
	event int
	data  string
	err   error
}

// String converts Packet object to a string format
func (packet *Packet) String() string {
	return fmt.Sprintf("%d%s", packet.event, packet.data)
}

// Byte converts Packet object to byte array format
func (packet *Packet) Byte() []byte {
	return []byte(packet.String())
}
