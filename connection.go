package sphere

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// NewConnection returns a new ws connection instance
func NewConnection(upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request) (*Connection, IError) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		return &Connection{guid.String(), 0, newChannelMap(), make(chan *Packet), make(chan struct{}), r, ws}, nil
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
	channels channelmap
	// buffered channel of outbound messages
	send chan *Packet
	// done channel
	done chan struct{}
	// http request
	request *http.Request
	// websocket connection
	*websocket.Conn
}

// queue is the connection message queue
func (conn *Connection) queue() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case packet, ok := <-conn.send:
			if !ok {
				conn.emit(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.emit(websocket.TextMessage, packet); err != nil {
				LogError(err)
				return
			}
		case <-conn.done:
			close(conn.done)
			return
		case <-ticker.C:
			if err := conn.emit(websocket.PingMessage, []byte{}); err != nil {
				LogError(err)
				return
			}
		}
	}
}

// write

// write writes a message with the given message type and payload.
func (conn *Connection) emit(mt int, payload interface{}) IError {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	switch msg := payload.(type) {
	case []byte:
		return conn.WriteMessage(mt, msg)
	case *Packet:
		if msg == nil {
			return ErrBadScheme
		}
		if !msg.Reply {
			conn.cid++
		}
		msg.Cid = conn.cid
		json, err := msg.ToJSON()
		if err != nil {
			return err
		}
		return conn.WriteMessage(websocket.TextMessage, json)
	}
	return nil
}

// subscribe to channel
func (conn *Connection) subscribe(channel *Channel) IError {
	if !conn.isSubscribed(channel) {
		conn.channels.Set(channel.Name(), channel)
	}
	if !channel.connections.Has(conn.id) {
		channel.connections.Set(conn.id, conn)
	}
	return nil
}

// unsubscribe from channel
func (conn *Connection) unsubscribe(channel *Channel) IError {
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

// close connection
func (conn *Connection) close() {
	conn.done <- struct{}{}
}
