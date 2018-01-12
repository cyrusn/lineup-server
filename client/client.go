package client

import (
	"log"
	"net/http"

	"github.com/cyrusn/lineup-system/hub"
	"github.com/gorilla/websocket"
)

// Client ...
type Client struct {
	Hub  *hub.Hub
	Conn *websocket.Conn
}

// ServeWS is a handler for connection
func ServeWS(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	Upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	c := &Client{
		Hub:  h,
		Conn: conn,
	}
	c.Hub.Register <- conn
	go c.Read()
}

// Read read the connection to process close, ping and pong messages
func (c *Client) Read() {
	defer func() {
		c.Hub.Unregister <- c.Conn
		c.Conn.Close()
	}()

	for {
		err := c.Conn.ReadJSON(&c.Hub.Message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
				c.Hub.Unregister <- c.Conn
			}
			break
		}
		c.Hub.ChanMessage <- c.Hub.Message
	}
}
