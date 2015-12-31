package sphere

// ChannelError indicates the error of the channel
type ChannelError int

const (
	// ErrorAlreadySubscribed indicates that this channel is in a subscribed state
	ErrorAlreadySubscribed ChannelError = iota
	// ErrorNotSubscribed indicates that this channel is in a unsubscribed state
	ErrorNotSubscribed
)

// ChannelErrorCode returns the string value of SphereError
var ChannelErrorCode = [...]string{
	"channel is already subscribed",
	"channel is not subscribed",
}

// Returns the error message
func (s ChannelError) String() string {
	return ChannelErrorCode[s]
}
