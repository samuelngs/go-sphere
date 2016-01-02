package sphere

import (
	"encoding/json"
	"strings"
)

// PacketType indicates the error of the channel
type PacketType int

const (
	// PacketTypeMessage denotes a regular text data message.
	PacketTypeMessage PacketType = iota
	// PacketTypeChannel denotes a channel request.
	PacketTypeChannel
	// PacketTypeSubscribe denotes a subscribe request.
	PacketTypeSubscribe
	// PacketTypeUnsubscribe denotes an unsubscribe request.
	PacketTypeUnsubscribe
	// PacketTypeSubscribed denotes a response of subscribe request.
	PacketTypeSubscribed
	// PacketTypeUnsubscribed denotes a response of unsubscribe request.
	PacketTypeUnsubscribed
	// PacketTypePing denotes an ping message.
	PacketTypePing
	// PacketTypePong denotes an pong message.
	PacketTypePong
	// PacketTypeUnknown denotes an pong message.
	PacketTypeUnknown
)

// PacketTypeCode returns the string value of SphereError
var PacketTypeCode = [...]string{
	"message",
	"channel",
	"subscribe",
	"unsubscribe",
	"subscribed",
	"unsubscribed",
	"ping",
	"pong",
	"unknown",
}

// Returns the error message
func (p PacketType) String() string {
	return PacketTypeCode[p]
}

// UnmarshalJSON to parse object from json string
func (p *PacketType) UnmarshalJSON(b []byte) (err error) {
	switch strings.Trim(string(b[:]), `"`) {
	case PacketTypeCode[0]:
		*p = PacketTypeMessage
	case PacketTypeCode[1]:
		*p = PacketTypeChannel
	case PacketTypeCode[2]:
		*p = PacketTypeSubscribe
	case PacketTypeCode[3]:
		*p = PacketTypeUnsubscribe
	case PacketTypeCode[4]:
		*p = PacketTypeSubscribed
	case PacketTypeCode[5]:
		*p = PacketTypeUnsubscribed
	case PacketTypeCode[6]:
		*p = PacketTypePing
	case PacketTypeCode[6]:
		*p = PacketTypePong
	default:
		*p = PacketTypeUnknown
	}
	return nil
}

// MarshalJSON to convert PacketType to json string
func (p *PacketType) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}
