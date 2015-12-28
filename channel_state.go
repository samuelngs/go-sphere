package sphere

// ChannelState indicates the state of the channel
type ChannelState int

const (
	// ChannelStateSubscribed indicates that this channel is in a subscribed state
	ChannelStateSubscribed ChannelState = iota
	// ChannelStateUnsubscribed indicates that this channel is in a unsubscribed state
	ChannelStateUnsubscribed
	// ChannelStatePending indicates that this channel is in a pending state
	ChannelStatePending
)

// ChannelStateCode returns the string value of ChannelState
var ChannelStateCode = [...]string{
	"subscribed",
	"unsubscribed",
	"pending",
}

// Returns the code id of channel state
func (s ChannelState) String() string {
	return ChannelStateCode[s]
}
