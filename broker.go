package sphere

import "errors"

const (
	// BrokerErrorOverrideOnSubscribe to warn use to override OnSubsribe function
	BrokerErrorOverrideOnSubscribe = "please override OnSubscribe"
	// BrokerErrorOverrideOnUnsubscribe to warn use to override OnUnsubscribe function
	BrokerErrorOverrideOnUnsubscribe = "please override OnUnsubscribe"
	// BrokerErrorOverrideOnPublish to warn use to override OnPublish function
	BrokerErrorOverrideOnPublish = "please override OnPublish"
	// BrokerErrorOverrideOnMessage to warn use to override OnMessage function
	BrokerErrorOverrideOnMessage = "please override OnMessage"
)

// Agent represents Broker instance
type Agent interface {
	OnSubscribe(*Channel) error            // => Broker OnSubscribe
	OnUnsubscribe(*Channel) error          // => Broker OnUnsubscribe
	OnPublish(*Channel, interface{}) error // => Broker OnPublish
	OnMessage(*Channel, interface{}) error // => Broker OnMessage
}

// ExtendBroker creates a broker instance
func ExtendBroker() *Broker {
	return &Broker{
		channels: make(map[string]*Channel),
	}
}

// Broker allows you to interact directly with Websocket internal data and pub/sub channels
type Broker struct {
	// List of channels
	channels map[string]*Channel
}

// OnSubscribe when websocket subscribes to a channel
func (broker *Broker) OnSubscribe(channel *Channel) error {
	return errors.New(BrokerErrorOverrideOnSubscribe)
}

// OnUnsubscribe when websocket unsubscribes from a channel
func (broker *Broker) OnUnsubscribe(channel *Channel) error {
	return errors.New(BrokerErrorOverrideOnUnsubscribe)
}

// OnPublish when websocket publishes data to a particular channel from the current broker
func (broker *Broker) OnPublish(channel *Channel, data interface{}) error {
	return errors.New(BrokerErrorOverrideOnPublish)
}

// OnMessage when websocket receive data from the broker subscriber
func (broker *Broker) OnMessage(channel *Channel, data interface{}) error {
	return errors.New(BrokerErrorOverrideOnMessage)
}
