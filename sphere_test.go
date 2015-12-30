package sphere

import (
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	address = "127.0.0.1:53320"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	s := NewSphere()
	r := gin.New()
	r.GET("/sync", func(c *gin.Context) {
		s.Handler(c.Writer, c.Request)
	})
	go r.Run(address)
}

func CreateConnection() (c *websocket.Conn, response *http.Response, err error) {
	u, err := url.Parse("ws://" + address + "/sync")
	if err != nil {
		return nil, nil, err
	}
	rawConn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return nil, nil, err
	}
	wsHeaders := http.Header{
		"Origin": {"http://" + address},
	}
	wsConn, resp, err := websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)
	return wsConn, resp, err
}

func TestSphereReady(t *testing.T) {
	count := 0
	ticker := time.NewTicker(time.Millisecond * 60)
	for {
		select {
		case <-ticker.C:
			s, err := net.Dial("tcp", address)
			defer s.Close()
			if err == nil || count > 3 {
				defer ticker.Stop()
				if err != nil {
					t.Fatal("Server could not get started")
				}
				return
			}
		}
		count++
	}
}

func TestSphereConnection(t *testing.T) {
	_, _, err := CreateConnection()
	if err != nil {
		t.Fatal(err.Error())
	}
}
