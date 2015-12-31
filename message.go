package sphere

// Message indicates the data of the message
type Message struct {
	Event string      `json:"event,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
