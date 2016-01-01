package sphere

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

// Packet indicates the data of the message
type Packet struct {
	Type      PacketType `json:"type"`
	Namespace string     `json:"namespace,omitempty"`
	Room      string     `json:"room,omitempty"`
	Cid       int        `json:"cid,omitempty"`
	Rid       int        `json:"rid,omitempty"`
	Error     error      `json:"error,omitempty"`
	Message   *Message   `json:"message,omitempty"`
	Machine   string     `json:"-"`
}

// ParsePacket returns Packet from bytes
func ParsePacket(data []byte) (*Packet, error) {
	var p *Packet
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, errors.New("packet format is invalid")
	}
	return p, nil
}

// Packet.toJSON returns json byte array from Packet
func (p *Packet) toJSON() ([]byte, error) {
	return json.Marshal(p)
}

// Packet.toBytes returns byte array from Packet
func (p *Packet) toBytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
