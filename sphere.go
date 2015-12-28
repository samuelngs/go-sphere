package sphere

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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
	upgrader = websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
	}
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
		connections: make(map[string]*Connection),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
	}
	go sphere.broker()
	go sphere.queue()
	return sphere
}

// Sphere represents an entire Websocket instance
type Sphere struct {
	// a broker agent
	agent Agent
	// list of active connections
	connections map[string]*Connection
	// register requests from the connection(s)
	register chan *Connection
	// unregister requests from the connection(s)
	unregister chan *Connection
}

// Broker handler for Sphere
func (sphere *Sphere) broker() {
}

// Queue handler for Sphere
func (sphere *Sphere) queue() {
	for {
		select {
		case conn := <-sphere.register:
			sphere.connections[conn.id] = conn
		case conn := <-sphere.unregister:
			delete(sphere.connections, conn.id)
			defer conn.Close()
		}
	}
}

// Handler handles and creates websocket connection
func (sphere *Sphere) Handler(w http.ResponseWriter, r *http.Request) {
	if conn, err := NewConnection(w, r); err == nil {
		sphere.register <- conn
		go conn.writePump()
		conn.readPump()
		defer func() {
			sphere.unregister <- conn
		}()
	}
}
