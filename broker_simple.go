package sphere

// NewSimpleBroker creates a new instance of SimpleBroker
func NewSimpleBroker() *SimpleBroker {
	return &SimpleBroker{
		ExtendBroker(),
	}
}

// SimpleBroker is a broker adapter built on Simple client
type SimpleBroker struct {
	*Broker
}

// OnSubscribe when websocket subscribes to a channel
func (broker *SimpleBroker) OnSubscribe(channel *Channel) error {
	return nil
}

// OnUnsubscribe when websocket unsubscribes from a channel
func (broker *SimpleBroker) OnUnsubscribe(channel *Channel) error {
	return nil
}

// OnPublish when websocket publishes data to a particular channel from the current broker
func (broker *SimpleBroker) OnPublish(channel *Channel, data *Packet) error {
	return nil
}

// OnMessage when websocket receive data from the broker subscriber
func (broker *SimpleBroker) OnMessage(channel *Channel, data *Packet) error {
	return nil
}
