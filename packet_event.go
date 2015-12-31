package sphere

// PacketEvent indicates the event of the Packet
type PacketEvent int

const (
	// PacketEventMessage denotes a text data message.
	PacketEventMessage PacketEvent = iota
	// PacketEventSubscribe denotes a subscribe request data message.
	PacketEventSubscribe
	// PacketEventUnsubscribe denotes an unsubscribe request data message.
	PacketEventUnsubscribe
)

// PacketEventCode returns the string value of PacketEvent
var PacketEventCode = [...]string{
	"message",
	"subscribe",
	"unsubscribe",
}

// Returns the code id of packet event
func (s PacketEvent) String() string {
	return PacketEventCode[s]
}
