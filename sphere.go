package sphere

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Sphere represents an entire Websocket instance
type Sphere struct {
	// a broker agent
	agent *Agent
	// list of active connections
	connections []*Connection
}

// Handler handles and creates websocket connection
func (sphere *Sphere) Handler(w http.ResponseWriter, r *http.Request) {
	if _, err := NewConnection(w, r); err == nil {
	} else {
	}
}
