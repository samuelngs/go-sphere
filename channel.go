package sphere

import "github.com/streamrail/concurrent-map"

// NewChannel creates new Channel instance
func NewChannel(name string) *Channel {
	return &Channel{name: name, state: ChannelStatePending, connections: cmap.New()}
}

// Channel let you subscribe to and watch for incoming data which is published on that channel by other clients or the server
type Channel struct {
	name        string
	state       ChannelState
	connections cmap.ConcurrentMap
}

// Name returns the name of the channel
func (channel *Channel) Name() string {
	return channel.name
}

// State returns the state of the channel
func (channel *Channel) State() ChannelState {
	return channel.state
}

// Connections returns a list of active user connections
func (channel *Channel) Connections() []*Connection {
	conns := make([]*Connection, 0, len(channel.connections))
	for item := range channel.connections.Iter() {
		conns = append(conns, item.Val.(*Connection))
	}
	return conns
}
