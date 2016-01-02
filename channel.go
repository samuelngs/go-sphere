package sphere

// NewChannel creates new Channel instance
func NewChannel(namespace string, room string) *Channel {
	return &Channel{namespace: namespace, room: room, state: ChannelStatePending, connections: NewConnectionMap()}
}

// Channel let you subscribe to and watch for incoming data which is published on that channel by other clients or the server
type Channel struct {
	namespace   string
	room        string
	state       ChannelState
	connections ConnectionMap
}

// Name returns the name of the channel
func (channel *Channel) Name() string {
	return channel.namespace + ":" + channel.room
}

// State returns the state of the channel
func (channel *Channel) State() ChannelState {
	return channel.state
}

// Connections returns a list of active user connections
func (channel *Channel) Connections() []*Connection {
	conns := make([]*Connection, 0, len(channel.connections))
	for item := range channel.connections.Iter() {
		conns = append(conns, item.Val)
	}
	return conns
}

// subscribe this channel
func (channel *Channel) subscribe(c *Connection) error {
	state := channel.isSubscribed(c)
	if !state {
		c.subscribe(channel)
	}
	return nil
}

// unsubscribe this channel
func (channel *Channel) unsubscribe(c *Connection) error {
	state := channel.isSubscribed(c)
	if state {
		c.unsubscribe(channel)
	}
	return nil
}

// isSubscribed checks if connection is in the connection list
func (channel *Channel) isSubscribed(c *Connection) bool {
	return channel.connections.Has(c.id)
}

// Emit sends message to current channel
func (channel *Channel) emit(mt int, payload []byte, c *Connection) error {
	l := channel.connections.Count()
	e := make(chan error, l)
	go func() {
		for item := range channel.connections.Iter() {
			conn := item.Val
			if conn != c {
				e <- conn.emit(mt, payload)
			} else {
				e <- nil
			}
		}
	}()
	for i := 0; i < l; i++ {
		<-e
	}
	return nil
}
