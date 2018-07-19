package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client store connecton and Action send to client
type Client struct {
	Conn *websocket.Conn
	Hub  *Hub
	Send chan []byte
}

func newClient(conn *websocket.Conn, h *Hub) *Client {
	return &Client{
		conn,
		h,
		make(chan []byte),
	}
}

func (c *Client) read() {
	defer func() {
		c.Hub.Unregister <- c
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()

		if websocket.IsUnexpectedCloseError(
			err,
			websocket.CloseGoingAway,
			websocket.CloseAbnormalClosure,
		) {
			log.Printf("error: %v", err)
			c.Hub.Unregister <- c
			break
		}

		c.Hub.onBroadcast(msg)
	}
}

func (c *Client) write() {
	for {
		select {
		case data, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (c *Client) close() {
	c.Conn.Close()
	close(c.Send)
}

// Run run the write function for client
func (c *Client) Run() {
	// go c.read()
	go c.write()
}
