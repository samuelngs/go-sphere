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

type simpleBrokerPubSub struct {
	receive chan *Packet
	done    chan bool
}

// OnSubscribe when websocket subscribes to a channel
func (broker *SimpleBroker) OnSubscribe(channel *Channel) error {
	c := make(chan error)
	go func() {
		if broker.store.Has(channel.Name()) {
			c <- nil
			return
		}
		// creates subscribe pubsub
		pubsub := &simpleBrokerPubSub{receive: make(chan *Packet), done: make(chan bool)}
		broker.store.Set(channel.Name(), pubsub)
		for {
			select {
			case p := <-pubsub.receive:
				broker.OnMessage(channel, p)
			case <-pubsub.done:
				return
			}
		}
	}()
	return <-c
}

// OnUnsubscribe when websocket unsubscribes from a channel
func (broker *SimpleBroker) OnUnsubscribe(channel *Channel) error {
	c := make(chan error)
	go func() {
		if !broker.store.Has(channel.Name()) {
			c <- nil
			return
		}
		if tmp, ok := broker.store.Get(channel.Name()); ok {
			if pubsub, ok := tmp.(*simpleBrokerPubSub); ok {
				pubsub.done <- true
				close(pubsub.receive)
				close(pubsub.done)
				broker.store.Remove(channel.Name())
			}
		}
		c <- nil
	}()
	return <-c
}

// OnPublish when websocket publishes data to a particular channel from the current broker
func (broker *SimpleBroker) OnPublish(channel *Channel, data *Packet) error {
	c := make(chan error)
	go func() {
		if broker.store.Has(channel.Name()) {
			if tmp, ok := broker.store.Get(channel.Name()); ok {
				if pubsub, ok := tmp.(*simpleBrokerPubSub); ok {
					pubsub.receive <- data
				}
			}
		}
		c <- nil
	}()
	return <-c
}

// OnMessage when websocket receive data from the broker subscriber
func (broker *SimpleBroker) OnMessage(channel *Channel, data *Packet) error {
	c := make(chan error)
	go func() {
		if json, err := data.toJSON(); err == nil {
			channel.emit(TextMessage, json, nil)
		}
		c <- nil
	}()
	return <-c
}
