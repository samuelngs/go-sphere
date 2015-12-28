package sphere

import (
	redis "gopkg.in/redis.v3"
)

var (
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
)

// NewRedisBroker creates a new instance of RedisBroker
func NewRedisBroker() *RedisBroker {
	return &RedisBroker{
		NewBroker(),
	}
}

// RedisBroker is a broker adapter built on Redis client
type RedisBroker struct {
	*Broker
}

// OnSubscribe when websocket subscribes to a channel
func (broker *RedisBroker) OnSubscribe() error {
	return nil
}

// OnUnsubscribe when websocket unsubscribes from a channel
func (broker *RedisBroker) OnUnsubscribe() error {
	return nil
}

// OnPublish when websocket publishes data to a particular channel from the current broker
func (broker *RedisBroker) OnPublish(channel *Channel, data interface{}) error {
	return nil
}
