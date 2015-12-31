package sphere

import redis "gopkg.in/redis.v3"

var (
	roption = redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	pubclient = redis.NewClient(&roption)
	subclient = redis.NewClient(&roption)
)

// NewRedisBroker creates a new instance of RedisBroker
func NewRedisBroker() *RedisBroker {
	return &RedisBroker{
		ExtendBroker(),
	}
}

// RedisBroker is a broker adapter built on Redis client
type RedisBroker struct {
	*Broker
}

// OnSubscribe when websocket subscribes to a channel
func (broker *RedisBroker) OnSubscribe(channel *Channel) error {
	c := make(chan error)
	go func() {
		// return if pubsub is already existed
		if broker.store.Has(channel.name) {
			c <- nil
			return
		}
		// creates subscribe pubsub
		pubsub, err := subclient.Subscribe(channel.name)
		if err == nil {
			broker.store.Set(channel.name, pubsub)
		}
		// close pubsub when process is done
		defer broker.OnUnsubscribe(channel)
		for {
			_, err := pubsub.ReceiveMessage()
			if err != nil {
				c <- err
				return
			}
		}
	}()
	return <-c
}

// OnUnsubscribe when websocket unsubscribes from a channel
func (broker *RedisBroker) OnUnsubscribe(channel *Channel) error {
	c := make(chan error)
	go func() {
		if tmp, ok := broker.store.Get(channel.name); ok {
			if pubsub, ok := tmp.(*redis.PubSub); ok {
				// close pubsub handler
				c <- pubsub.Close()
				// remove pubsub from store
				broker.store.Remove(channel.name)
			}
		} else {
			c <- nil
		}
	}()
	return <-c
}

// OnPublish when websocket publishes data to a particular channel from the current broker
func (broker *RedisBroker) OnPublish(channel *Channel, data interface{}) error {
	return nil
}

// OnMessage when websocket receive data from the broker subscriber
func (broker *RedisBroker) OnMessage(channel *Channel, data interface{}) error {
	return nil
}
