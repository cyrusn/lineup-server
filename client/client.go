package client

import (
	"log"
	"net/http"

	"github.com/cyrusn/lineup-system/chatroom"
	"github.com/cyrusn/lineup-system/hub"
	"github.com/cyrusn/lineup-system/schedule"
	"github.com/gorilla/websocket"
)

// Client ...
type Client struct {
	Hub          *hub.Hub
	Conn         *websocket.Conn
	MapSchedules schedule.MapSchedules
	Message      chatroom.Message
}

// ServeWS is a handler for connection
func ServeWS(h *hub.Hub, s schedule.MapSchedules, w http.ResponseWriter, r *http.Request) {
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
		Hub:          h,
		Conn:         conn,
		MapSchedules: s,
		Message:      chatroom.Message{},
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
		go c.ReadMessage()
	}
}

// ReadSchedule read schedule update from client
func (c *Client) ReadSchedule() {
	err := c.Conn.ReadJSON(c.MapSchedules)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
		return
	}
	c.Hub.ChanMapScheudle <- c.MapSchedules
}

// ReadMessage read message update from client
func (c *Client) ReadMessage() {
	err := c.Conn.ReadJSON(c.Message)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
		return
	}
	c.Hub.ChanMessage <- c.Message
}

// BroadcastSchedule broadcast schedule
func (c *Client) BroadcastSchedule() {
	c.Hub.ChanMapScheudle <- c.MapSchedules
}
