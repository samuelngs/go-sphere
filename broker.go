package sphere

// Broker allows you to interact directly with Websocket internal data and pub/sub channels
type Broker struct {
	// The broker's id
	ID string
	// List of channels
	Channels map[string]*Channel
}
