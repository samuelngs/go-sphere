package sphere

// IChannels is the interface for ChannelModel
type IChannels interface {
	Name() string
	Subscribe(uri string, connection *Connection) bool
	Disconnect(uri string, connection *Connection)
	Receive(event string, message string) error
}

// ExtendChannelModel lets developer create a IChannals compatible struct
func ExtendChannelModel(name string) *ChannelModel {
	return &ChannelModel{name}
}

// ChannelModel is for user to define channel events and actions
type ChannelModel struct {
	name string
}

// Name to return name of the channel
func (m *ChannelModel) Name() string {
	return m.name
}

// Subscribe decides whether accept the connection into channel or not, return true => accept, false => reject
func (m *ChannelModel) Subscribe(uri string, connection *Connection) bool {
	return false
}

// Disconnect defines the action when user disconnect from channel
func (m *ChannelModel) Disconnect(uri string, connection *Connection) {
}

// Receive defines the action when websocket server receive message from user in this channel
func (m *ChannelModel) Receive(event string, message string) error {
	return nil
}
