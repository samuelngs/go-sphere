package sphere

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	// PacketEventMessage denotes a text data message.
	PacketEventMessage = 0
	// PacketEventSubscribe denotes a subscribe request data message.
	PacketEventSubscribe = 6
	// PacketEventUnsubscribe denotes an unsubscribe request data message.
	PacketEventUnsubscribe = 7
)

// NewPacket creates new packet instance
func NewPacket(event int, data string) *Packet {
	return &Packet{event, data}
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
	if i != PacketEventMessage && i != PacketEventSubscribe && i != PacketEventUnsubscribe {
		return nil, errors.New("packet is invalid")
	}
	return NewPacket(i, str[1:len(str)]), nil
}

// Packet indicates the data of the message
type Packet struct {
	event int
	data  string
}

// String converts Packet object into a string format
func (packet *Packet) String() string {
	return fmt.Sprintf("%d%s", packet.event, packet.data)
}
