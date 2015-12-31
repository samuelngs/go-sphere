package sphere

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	address = func() string {
		l, _ := net.Listen("tcp", ":0")
		defer l.Close()
		return fmt.Sprintf("127.0.0.1:%d", l.Addr().(*net.TCPAddr).Port)
	}()
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	a := NewRedisBroker()
	s, r := NewSphere(a), gin.New()
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
	return websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)
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
	c, _, err := CreateConnection()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer c.Close()
}

func TestSphereSendMessage(t *testing.T) {
	c, _, err := CreateConnection()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer c.Close()
	p := &Packet{Success: true, Type: PacketTypePing}
	res, err := p.toJSON()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := c.WriteMessage(TextMessage, res); err != nil {
		t.Fatal(err.Error())
	}
}

func TestSphereMessagePingPong(t *testing.T) {
	done := make(chan error)
	c, _, err := CreateConnection()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer c.Close()
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				done <- err
				return
			}
			done <- nil
			fmt.Println("ahah", string(msg[:]))
			return
		}
	}()
	p := &Packet{Success: true, Type: PacketTypePing}
	res, err := p.toJSON()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := c.WriteMessage(TextMessage, res); err != nil {
		t.Fatal(err.Error())
	}
	<-done
}
