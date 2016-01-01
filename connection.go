package sphere

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/streamrail/concurrent-map"
)

// NewConnection returns a new ws connection instance
func NewConnection(w http.ResponseWriter, r *http.Request) (*Connection, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		conn := &Connection{guid.String(), 0, cmap.New(), make(chan []byte), make(chan *Packet), r, ws}
		go conn.queue()
		return conn, nil
	}
	return nil, err
}

// Connection allows you to interact with backend and other client sockets in realtime
type Connection struct {
	// the id of the connection
	id string
	// cid
	cid int
	// list of channels that this connection has been subscribed
	channels cmap.ConcurrentMap
	// buffered channel of outbound messages
	send chan []byte
	// buffered channel of inbound messages
	receive chan *Packet
	// http request
	request *http.Request
	// websocket connection
	*websocket.Conn
}

func (conn *Connection) queue() {
	for {
		select {
		case <-conn.receive:
		}
	}
}

// write writes a message with the given message type and payload.
func (conn *Connection) emit(mt int, payload []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.WriteMessage(mt, payload)
}

// subscribe to channel
func (conn *Connection) subscribe(channel *Channel) error {
	if !conn.isSubscribed(channel) {
		conn.channels.Set(channel.Name(), channel)
	}
	if !channel.connections.Has(conn.id) {
		conn.channels.Set(conn.id, conn)
	}
	return nil
}

// unsubscribe from channel
func (conn *Connection) unsubscribe(channel *Channel) error {
	if conn.isSubscribed(channel) {
		conn.channels.Remove(channel.Name())
	}
	if channel.connections.Has(conn.id) {
		conn.channels.Remove(conn.id)
	}
	return nil
}

// isSubscribed checks if channel is subscribed
func (conn *Connection) isSubscribed(channel *Channel) bool {
	return conn.channels.Has(channel.Name())
}
