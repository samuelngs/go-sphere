package sphere

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	"github.com/streamrail/concurrent-map"
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
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1
	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2
	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8
	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9
	// PongMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
	// Message Types
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

// NewSphere creates a new instance of Sphere
func NewSphere(brokers ...Agent) *Sphere {
	// declare agent
	var broker Agent
	// set declared agent if parameter exists
	for _, i := range brokers {
		broker = i
		break
	}
	if broker == nil {
		broker = NewSimpleBroker()
	}
	// creates sphere instance
	sphere := &Sphere{
		agent:       broker,
		connections: cmap.New(),
		channels:    cmap.New(),
		models:      cmap.New(),
	}
	return sphere
}

// Sphere represents an entire Websocket instance
type Sphere struct {
	// a broker agent
	agent Agent
	// list of active connections
	connections cmap.ConcurrentMap
	// list of channels
	channels cmap.ConcurrentMap
	// list of models
	models cmap.ConcurrentMap
}

// Handler handles and creates websocket connection
func (sphere *Sphere) Handler(w http.ResponseWriter, r *http.Request) {
	if conn, err := NewConnection(w, r); err == nil {
		sphere.connections.Set(conn.id, conn)
		go sphere.write(conn)
		sphere.read(conn)
		defer sphere.connections.Remove(conn.id)
	}
}

// Models channel models
func (sphere *Sphere) Models(models ...IChannels) {
	for _, model := range models {
		if !sphere.models.Has(model.Namespace()) {
			sphere.models.Set(model.Namespace(), model)
		} else {
			panic(fmt.Sprintf("model \"%s\" is already existed", model.Namespace()))
		}
	}
}

func (sphere *Sphere) write(conn *Connection) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case msg, ok := <-conn.send:
			if !ok {
				conn.emit(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.emit(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := conn.emit(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (sphere *Sphere) read(conn *Connection) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if msg != nil {
			go sphere.process(conn, msg)
		}
	}
}

func (sphere *Sphere) process(conn *Connection, msg []byte) {
	p, err := ParsePacket(msg)
	if err != nil {
		fmt.Printf("Error: %v - %v", err.Error(), string(msg[:]))
		return
	}
	if p != nil {
		switch p.Type {
		case PacketTypeChannel:
			if p.Namespace != "" && p.Room != "" {
				if sphere.models.Has(p.Namespace) {
					p.Machine = sphere.agent.ID()
					sphere.publish(p)
				} else {
					p.Error = ErrUnsupportedNamespace
					if json, err := p.toJSON(); err == nil {
						conn.send <- json
					}
				}
			} else {
				p.Error = ErrBadScheme
				if json, err := p.toJSON(); err == nil {
					conn.send <- json
				}
			}
		case PacketTypeSubscribe:
			if p.Namespace != "" && p.Room != "" {
				if sphere.models.Has(p.Namespace) {
					sphere.subscribe(p.Namespace, p.Room, conn)
				} else {
					p.Error = ErrUnsupportedNamespace
					if json, err := p.toJSON(); err == nil {
						conn.send <- json
					}
				}
			} else {
				p.Error = ErrBadScheme
				if json, err := p.toJSON(); err == nil {
					conn.send <- json
				}
			}
		case PacketTypeUnsubscribe:
			if p.Namespace != "" && p.Room != "" {
				if sphere.models.Has(p.Namespace) {
					sphere.unsubscribe(p.Namespace, p.Room, conn)
				} else {
					p.Error = ErrUnsupportedNamespace
					if json, err := p.toJSON(); err == nil {
						conn.send <- json
					}
				}
			} else {
				p.Error = ErrBadScheme
				if json, err := p.toJSON(); err == nil {
					conn.send <- json
				}
			}
		case PacketTypePing:
			p.Type = PacketTypePong
			if json, err := p.toJSON(); err == nil {
				conn.send <- json
			}
		case PacketTypeMessage:
			conn.receive <- p
		}
	}
}

func (sphere *Sphere) channel(namespace string, room string, autoCreateOpts ...bool) *Channel {
	c := make(chan *Channel)
	autoCreateOpt := false
	for _, opt := range autoCreateOpts {
		autoCreateOpt = opt
		break
	}
	name := sphere.agent.ChannelName(namespace, room)
	go func() {
		if tmp, ok := sphere.channels.Get(name); ok {
			c <- tmp.(*Channel)
		} else {
			if autoCreateOpt {
				channel := NewChannel(name)
				sphere.channels.Set(name, channel)
				c <- channel
			} else {
				c <- nil
			}
		}
	}()
	return <-c
}

func (sphere *Sphere) subscribe(namespace string, room string, conn *Connection) error {
	if channel := sphere.channel(namespace, room, true); channel != nil {
		if !sphere.agent.IsSubscribed(channel.name) {
			sphere.agent.OnSubscribe(channel)
		}
	}
	return nil
}

func (sphere *Sphere) unsubscribe(namespace string, room string, conn *Connection) error {
	if channel := sphere.channel(namespace, room, true); channel != nil {
		if sphere.agent.IsSubscribed(channel.name) {
			sphere.agent.OnUnsubscribe(channel)
		}
	}
	return nil
}

func (sphere *Sphere) publish(p *Packet) error {
	if channel := sphere.channel(p.Namespace, p.Room); channel != nil {
		sphere.agent.OnPublish(channel, p)
	}
	return nil
}
