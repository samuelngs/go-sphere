package sphere

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// NewConnection returns a new ws connection instance
func NewConnection(w http.ResponseWriter, r *http.Request) (*Connection, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		conn := &Connection{guid.String(), 0, NewChannelMap(), make(chan []byte), make(chan *Packet), r, ws}
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
	channels ChannelMap
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
		case data := <-conn.send:
			conn.emit(TextMessage, data)
		}
	}
}

// write writes a message with the given message type and payload.
func (conn *Connection) emit(mt int, payload interface{}, responses ...bool) error {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	switch msg := payload.(type) {
	case []byte:
		conn.cid++
		return conn.WriteMessage(mt, msg)
	case *Packet:
		response := false
		for _, r := range responses {
			response = r
			break
		}
		if msg == nil {
			return ErrBadScheme
		}
		conn.cid++
		if !response {
			msg.Cid = conn.cid
		}
		json, err := msg.toJSON()
		if err != nil {
			return err
		}
		return conn.WriteMessage(TextMessage, json)
	}
	return nil
}

// subscribe to channel
func (conn *Connection) subscribe(channel *Channel) error {
	if !conn.isSubscribed(channel) {
		conn.channels.Set(channel.Name(), channel)
	}
	if !channel.connections.Has(conn.id) {
		channel.connections.Set(conn.id, conn)
	}
	return nil
}

// unsubscribe from channel
func (conn *Connection) unsubscribe(channel *Channel) error {
	if conn.isSubscribed(channel) {
		conn.channels.Remove(channel.Name())
	}
	if channel.connections.Has(conn.id) {
		channel.connections.Remove(conn.id)
	}
	return nil
}

// isSubscribed checks if channel is subscribed
func (conn *Connection) isSubscribed(channel *Channel) bool {
	return conn.channels.Has(channel.Name())
}
