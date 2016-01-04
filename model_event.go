package sphere

// IEvents represents EventModel
type IEvents interface {
	Namespace() string
	Receive(string, string) (string, IError)
}

// ExtendEventModel lets developer create a IEvents compatible struct
func ExtendEventModel(namespace string) *EventModel {
	return &EventModel{namespace}
}

// EventModel is for user to define channel events and actions
type EventModel struct {
	namespace string
}

// Namespace to return name of the channel
func (m *EventModel) Namespace() string {
	return m.namespace
}

// Receive defines the action when websocket server receive message from user in this channel
func (m *EventModel) Receive(event string, message string) (string, IError) {
	return "", nil
}
