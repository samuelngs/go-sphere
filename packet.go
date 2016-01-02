package sphere

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

// Packet indicates the data of the message
type Packet struct {
	Type      PacketType
	Namespace string
	Room      string
	Cid       int
	Error     error
	Message   *Message
	Machine   string
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

// MarshalJSON handler
func (p *Packet) MarshalJSON() ([]byte, error) {
	var err string
	if p.Error != nil {
		err = p.Error.Error()
	}
	return json.Marshal(&struct {
		Type      PacketType `json:"type"`
		Namespace string     `json:"namespace,omitempty"`
		Room      string     `json:"room,omitempty"`
		Cid       int        `json:"cid"`
		Error     string     `json:"error,omitempty"`
		Message   *Message   `json:"message,omitempty"`
		Machine   string     `json:"-"`
	}{p.Type, p.Namespace, p.Room, p.Cid, err, p.Message, p.Machine})
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

func (p *Packet) String() string {
	if json, err := p.toJSON(); err == nil {
		return string(json[:])
	}
	return ""
}

// Response return response packet
func (p *Packet) Response() *Packet {
	r := *p
	switch r.Type {
	case PacketTypeSubscribe:
		r.Type = PacketTypeSubscribed
	case PacketTypeUnsubscribe:
		r.Type = PacketTypeUnsubscribed
	case PacketTypePing:
		r.Type = PacketTypePong
	}
	return &r
}

// SetError set error message
func (p *Packet) SetError(err error) *Packet {
	p.Error = err
	return p
}
