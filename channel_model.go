package sphere

// IChannels is the interface for ChannelModel
type IChannels interface {
	Namespace() string
	Subscribe(string, *Connection) bool
	Disconnect(string, *Connection) bool
	Receive(string, string) (interface{}, error)
}

// ExtendChannelModel lets developer create a IChannals compatible struct
func ExtendChannelModel(namespace string) *ChannelModel {
	return &ChannelModel{namespace}
}

// ChannelModel is for user to define channel events and actions
type ChannelModel struct {
	namespace string
}

// Namespace to return name of the channel
func (m *ChannelModel) Namespace() string {
	return m.namespace
}

// Subscribe decides whether accept the connection into channel or not, return true => accept, false => reject
func (m *ChannelModel) Subscribe(room string, connection *Connection) bool {
	return false
}

// Disconnect defines the action when user disconnect from channel
func (m *ChannelModel) Disconnect(room string, connection *Connection) bool {
	return true
}

// Receive defines the action when websocket server receive message from user in this channel
func (m *ChannelModel) Receive(event string, message string) (interface{}, error) {
	return nil, nil
}
