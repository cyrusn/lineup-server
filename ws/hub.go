package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Hub is a hub store all clients
type Hub struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	Clients    map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Broadcast:  make(chan []byte),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.onConnect(c)
		case c := <-h.Unregister:
			h.onDisconnect(c)
		case m := <-h.Broadcast:
			h.onBroadcast(m)
		}
	}
}

func (h *Hub) onConnect(c *Client) {
	h.Clients[c] = true
}

func (h *Hub) onDisconnect(c *Client) {
	_, ok := h.Clients[c]
	if ok {
		delete(h.Clients, c)
		close(c.Send)
	}
}

func (h *Hub) onBroadcast(msg []byte) {
	for c, ok := range h.Clients {
		if ok {
			c.Send <- msg
		}
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not upgrade", http.StatusInternalServerError)
		return
	}

	c := newClient(conn, h)
	h.Register <- c
	c.Run()
}
