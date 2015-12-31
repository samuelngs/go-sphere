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
		conn := &Connection{guid.String(), []*Channel{}, make(chan []byte), make(chan *Packet), r, ws}
		go conn.queue()
		return conn, nil
	}
	return nil, err
}

// Connection allows you to interact with backend and other client sockets in realtime
type Connection struct {
	// the id of the connection
	id string
	// list of channels that this connection has been subscribed
	channels []*Channel
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
		case p := <-conn.receive:
			p.Type = PacketTypePong
			res, err := p.toJSON()
			if err != nil {
				return
			}
			conn.WriteMessage(TextMessage, res)
		}
	}
}

// write writes a message with the given message type and payload.
func (conn *Connection) emit(mt int, payload []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.WriteMessage(mt, payload)
}
