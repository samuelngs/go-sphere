package sphere

import "errors"

const (
	// BrokerErrorOverrideOnSubscribe to warn use to override OnSubsribe function
	BrokerErrorOverrideOnSubscribe = "please override OnSubscribe"
	// BrokerErrorOverrideOnUnsubscribe to warn use to override OnUnsubscribe function
	BrokerErrorOverrideOnUnsubscribe = "please override OnUnsubscribe"
	// BrokerErrorOverrideOnPublish to warn use to override OnPublish function
	BrokerErrorOverrideOnPublish = "please override OnPublish"
)

// Broker allows you to interact directly with Websocket internal data and pub/sub channels
type Broker struct {
	// The broker's id
	id string
	// List of channels
	channels map[string]*Channel
}

// OnSubscribe when websocket subscribes to a channel
func (broker *Broker) OnSubscribe() error {
	return errors.New(BrokerErrorOverrideOnSubscribe)
}

// OnUnsubscribe when websocket unsubscribes from a channel
func (broker *Broker) OnUnsubscribe() error {
	return errors.New(BrokerErrorOverrideOnUnsubscribe)
}

// OnPublish when websocket publishes data to a particular channel from the current broker
func (broker *Broker) OnPublish(channel *Channel, data interface{}) error {
	return errors.New(BrokerErrorOverrideOnPublish)
}
