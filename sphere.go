package sphere

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

const (
	// Read buffer size for websocket upgrader
	readBufferSize = 1024
	// Write buffer size for websocker upgrader
	writeBufferSize = 1024
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	// Websocket Upgrader
	upgrader = websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
	}
	// Guid to generate globally unique id
	guid = xid.New()
)

// Default creates a new instance of Sphere
func Default(brokers ...IBroker) *Sphere {
	// declare agent
	var broker IBroker
	// set declared agent if parameter exists
	for _, i := range brokers {
		broker = i
		break
	}
	if broker == nil {
		broker = DefaultSimpleBroker()
	}
	// creates sphere instance
	sphere := &Sphere{
		broker:      broker,
		connections: newConnectionMap(),
		channels:    newChannelMap(),
		models:      newChannelModelMap(),
		events:      newEventModelMap(),
	}
	return sphere
}

// Sphere represents an entire Websocket instance
type Sphere struct {
	// a broker agent
	broker IBroker
	// list of active connections
	connections connectionmap
	// list of channels
	channels channelmap
	// list of models
	models channelmodelmap
	// list of events
	events eventmodelmap
}

// Handler handles and creates websocket connection
func (sphere *Sphere) Handler(w http.ResponseWriter, r *http.Request) IError {
	if conn, err := NewConnection(w, r); err == nil {
		sphere.connections.Set(conn.id, conn)
		// run connection queue
		go conn.queue()
		// action after connection disconnected
		defer func() {
			// unsubscribe all channels
			for item := range conn.channels.Iter() {
				channel := item.Val
				sphere.unsubscribe(channel.namespace, channel.room, conn)
			}
			// close all send and receive buffers
			conn.close()
			// remove connection from sphere after disconnect
			sphere.connections.Remove(conn.id)
		}()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			if msg != nil {
				go sphere.process(conn, msg)
			}
		}
	} else {
		return err
	}
	return nil
}

// Models load channel or event models
func (sphere *Sphere) Models(models ...interface{}) {
	for _, item := range models {
		switch model := item.(type) {
		case IChannels:
			if !sphere.models.Has(model.Namespace()) {
				sphere.models.Set(model.Namespace(), model)
			}
		case IEvents:
			if !sphere.events.Has(model.Namespace()) {
				sphere.events.Set(model.Namespace(), model)
			}
		}
	}
}

// process parses and processes received message
func (sphere *Sphere) process(conn *Connection, msg []byte) {
	// convert received bytes to Packet object
	p, err := ParsePacket(msg)
	if err != nil {
		LogError(err)
		return
	}
	switch p.Type {
	case PacketTypeChannel:
		if p.Namespace != "" && p.Room != "" {
			// publish message to broker if it is a channel event / message
			p.Machine = sphere.broker.ID()
			sphere.publish(p, conn)
		} else {
			// if namespace or room is not provided, return error message
			p.Error = ErrBadScheme
			conn.send <- p
		}
	case PacketTypeSubscribe:
		if p.Namespace != "" && p.Room != "" {
			// subscribe connection to channel
			err := sphere.subscribe(p.Namespace, p.Room, conn)
			r := p.Response()
			r.SetError(err)
			// return success or failure message to user
			conn.send <- r
		} else {
			// if namespace or room is not provided, return error message
			p.Error = ErrBadScheme
			conn.send <- p
		}
	case PacketTypeUnsubscribe:
		if p.Namespace != "" && p.Room != "" {
			// unsubscribe connection from channel
			if sphere.models.Has(p.Namespace) {
				sphere.unsubscribe(p.Namespace, p.Room, conn)
			} else {
				// if namespace model does not existed, return error
				p.Error = ErrNotSupported
				conn.send <- p
			}
		} else {
			// if namespace or room is not provided, return error message
			p.Error = ErrBadScheme
			conn.send <- p
		}
	case PacketTypeMessage:
		if p.Namespace != "" {
			// receive event message
			sphere.receive(p, conn)
		} else {
			p.Error = ErrBadScheme
			conn.send <- p
		}
	case PacketTypePing:
		// ping-pong
		r := p.Response()
		conn.send <- r
	}
}

// channel returns Channel object, channel will be automatually created when autoCreateOpts is true
func (sphere *Sphere) channel(namespace string, room string, autoCreateOpts ...bool) *Channel {
	c := make(chan *Channel)
	autoCreateOpt := false
	for _, opt := range autoCreateOpts {
		autoCreateOpt = opt
		break
	}
	name := sphere.broker.ChannelName(namespace, room)
	go func() {
		if tmp, ok := sphere.channels.Get(name); ok {
			c <- tmp
		} else {
			if autoCreateOpt {
				channel := NewChannel(namespace, room)
				sphere.channels.Set(name, channel)
				c <- channel
			} else {
				c <- nil
			}
		}
	}()
	return <-c
}

// subscribe trigger Broker OnSubscribe action and put connection into channel connections list
func (sphere *Sphere) subscribe(namespace string, room string, conn *Connection) IError {
	var model IChannels
	if !sphere.models.Has(namespace) {
		return ErrNotSupported
	}
	if tmp, ok := sphere.models.Get(namespace); ok {
		model = tmp
	} else {
		return ErrNotSupported
	}
	if accept, err := model.Subscribe(room, conn); !accept && err == nil {
		return ErrUnauthorized
	} else if !accept && err != nil {
		return err
	}
	channel := sphere.channel(namespace, room, true)
	if channel == nil {
		return ErrNotFound
	}
	if err := channel.subscribe(conn); err != nil {
		return err
	}
	if !sphere.broker.IsSubscribed(channel.namespace, channel.room) {
		c := make(chan IError)
		go sphere.broker.OnSubscribe(channel, c)
		return <-c
	}
	return nil
}

// unsubscribe trigger Broker OnUnsubscribe action and remove connection from channel connections list
func (sphere *Sphere) unsubscribe(namespace string, room string, conn *Connection) IError {
	var model IChannels
	if !sphere.models.Has(namespace) {
		return ErrNotSupported
	}
	if tmp, ok := sphere.models.Get(namespace); ok {
		model = tmp
	} else {
		return ErrNotSupported
	}
	if err := model.Disconnect(room, conn); err != nil {
		return err
	}
	channel := sphere.channel(namespace, room, false)
	if channel == nil {
		return ErrNotFound
	}
	err := channel.unsubscribe(conn)
	if err == nil && channel.connections.Count() == 0 {
		if sphere.broker.IsSubscribed(channel.namespace, channel.room) {
			c := make(chan IError)
			go sphere.broker.OnUnsubscribe(channel, c)
			return <-c
		}
	}
	return err
}

// publish trigger Broker OnPublish action, send message to user from broker
func (sphere *Sphere) publish(p *Packet, conn *Connection) IError {
	var model IChannels
	if !sphere.models.Has(p.Namespace) {
		return ErrNotSupported
	}
	if tmp, ok := sphere.models.Get(p.Namespace); ok {
		model = tmp
	} else {
		return ErrNotSupported
	}
	channel := sphere.channel(p.Namespace, p.Room)
	if channel == nil {
		return ErrBadStatus
	}
	if !channel.isSubscribed(conn) {
		return ErrNotSubscribed
	}
	msg := p.Message
	if msg == nil || msg.Event == "" {
		return ErrBadScheme
	}
	res, err := model.Receive(msg.Event, msg.Data)
	if err != nil {
		return err
	}
	d := p.Response()
	if res != "" {
		d.Message.Data = res
	}
	if sphere.broker.IsSubscribed(channel.namespace, channel.room) {
		return sphere.broker.OnPublish(channel, d)
	}
	return ErrServerErrors
}

// receive message and event handler
func (sphere *Sphere) receive(p *Packet, conn *Connection) IError {
	var model IEvents
	if !sphere.events.Has(p.Namespace) {
		return ErrNotSupported
	}
	if tmp, ok := sphere.events.Get(p.Namespace); ok {
		model = tmp
	} else {
		return ErrNotSupported
	}
	msg := p.Message
	if msg == nil || msg.Event == "" {
		return ErrBadScheme
	}
	res, err := model.Receive(msg.Event, msg.Data)
	if err != nil {
		return err
	}
	d := p.Response()
	if res != "" {
		d.Message.Data = res
	}
	conn.send <- d
	return nil
}
