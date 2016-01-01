package sphere

// Message indicates the data of the message
type Message struct {
	Event string `json:"event,omitempty"`
	Data  string `json:"data,omitempty"`
}
