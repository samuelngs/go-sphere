package sphere

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// NewConnection returns a new ws connection instance
func NewConnection(w http.ResponseWriter, r *http.Request) (*Connection, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		return &Connection{[]*Channel{}, r, conn}, nil
	}
	return nil, err
}

// Connection allows you to interact with backend and other client sockets in realtime
type Connection struct {
	// list of channels that this connection has been subscribed
	channels []*Channel
	// http header
	*http.Request
	// websocket connection
	*websocket.Conn
}
