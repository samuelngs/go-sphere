package sphere

// ChannelEvent indicates the state of the channel
type ChannelEvent int

const (
	// ChannelEventSubscribe indicates that this channel is in a subscribed state
	ChannelEventSubscribe ChannelEvent = iota
	// ChannelEventUnsubscribe indicates that this channel is in a unsubscribed state
	ChannelEventUnsubscribe
	// ChannelEventSubscribeFail indicates that this channel is in a pending state
	ChannelEventSubscribeFail
)

// ChannelEventCode returns the string value of ChannelEvent
var ChannelEventCode = [...]string{
	"subscribe",
	"unsubscribe",
	"subscribeFail",
}

// Returns the code id of channel state
func (s ChannelEvent) String() string {
	return ChannelEventCode[s]
}
