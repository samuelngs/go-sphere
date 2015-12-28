package sphere

// Channel let you subscribe to and watch for incoming data which is published on that channel by other clients or the server
type Channel struct {
	name        string
	state       ChannelState
	connections map[string]*Connection
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
	for _, conn := range channel.connections {
		conns = append(conns, conn)
	}
	return conns
}
