package sphere

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

var (
	guid = xid.New()
)

// NewConnection returns a new ws connection instance
func NewConnection(w http.ResponseWriter, r *http.Request) (*Connection, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		return &Connection{guid.String(), []*Channel{}, make(chan []byte), r, conn}, nil
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
	// http request
	request *http.Request
	// websocket connection
	*websocket.Conn
}

// write writes a message with the given message type and payload.
func (conn *Connection) emit(mt int, payload []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.WriteMessage(mt, payload)
}

// writePump pumps messages from the sphere to the websocket connection.
func (conn *Connection) writePump() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case msg, ok := <-conn.send:
			if !ok {
				return conn.emit(websocket.CloseMessage, []byte{})
			}
			if err := conn.emit(websocket.TextMessage, msg); err != nil {
				return err
			}
		case <-ticker.C:
			if err := conn.emit(websocket.PingMessage, []byte{}); err != nil {
				return err
			}
		}
	}
}

func (conn *Connection) readPump() {
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPingHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
		}
		if msg != nil {
		}
	}
}
