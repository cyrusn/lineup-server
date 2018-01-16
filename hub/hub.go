package hub

import (
	"github.com/cyrusn/lineup-system/chatroom"
	"github.com/cyrusn/lineup-system/schedule"
	"github.com/gorilla/websocket"
)

// Hub is a hub store everything needs for create a websocket server
type Hub struct {
	Register        chan *websocket.Conn
	Unregister      chan *websocket.Conn
	ChanMapScheudle chan schedule.MapSchedules
	ChanMessage     chan chatroom.Message
	Clients         map[*websocket.Conn]bool
}

// New create a new hub
func New() *Hub {
	return &Hub{
		Register:        make(chan *websocket.Conn),
		Unregister:      make(chan *websocket.Conn),
		ChanMapScheudle: make(chan schedule.MapSchedules),
		ChanMessage:     make(chan chatroom.Message),
		Clients:         make(map[*websocket.Conn]bool),
	}
}

// Run select chan type to do relative action
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Clients[c] = true
			break
		case c := <-h.Unregister:
			if _, ok := h.Clients[c]; ok {
				c.Close()
				delete(h.Clients, c)
			}
			break
		case c := <-h.ChanMapScheudle:
			for client := range h.Clients {
				client.WriteJSON(c)
			}
			break
		case c := <-h.ChanMessage:
			for client := range h.Clients {
				client.WriteJSON(c)
			}
			break
		}
	}
}
