package sphere

import (
	redis "gopkg.in/redis.v3"
)

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
		make(map[*Channel]*redis.PubSub),
		ExtendBroker(),
	}
}

// RedisBroker is a broker adapter built on Redis client
type RedisBroker struct {
	store map[*Channel]*redis.PubSub
	*Broker
}

// OnSubscribe when websocket subscribes to a channel
func (broker *RedisBroker) OnSubscribe(channel *Channel) error {
	c := make(chan error)
	go func() {
		pubsub, err := subclient.Subscribe(channel.name)
		if err == nil {
			broker.store[channel] = pubsub
		}
		c <- err
	}()
	return <-c
}

// OnUnsubscribe when websocket unsubscribes from a channel
func (broker *RedisBroker) OnUnsubscribe(channel *Channel) error {
	c := make(chan error)
	go func() {
		if pubsub := broker.store[channel]; pubsub != nil {
			// close pubsub handler
			c <- pubsub.Close()
			// remove pubsub from store
			defer delete(broker.store, channel)
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
