package sphere

// IChannels is the interface for ChannelModel
type IChannels interface {
	Namespace() string
	Subscribe(string, *Message, *Connection) (bool, IError)
	Disconnect(string, *Connection) IError
	Receive(string, string) (string, IError)
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
func (m *ChannelModel) Subscribe(room string, message *Message, connection *Connection) (bool, IError) {
	return false, nil
}

// Disconnect defines the action when user disconnect from channel
func (m *ChannelModel) Disconnect(room string, connection *Connection) IError {
	return nil
}

// Receive defines the action when websocket server receive message from user in this channel
func (m *ChannelModel) Receive(event string, message string) (string, IError) {
	return "", nil
}
