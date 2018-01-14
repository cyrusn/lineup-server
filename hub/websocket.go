package hub

import (
	"log"
	"net/http"

	"github.com/cyrusn/lineup-system/schedule"
	"github.com/gorilla/websocket"
)

// Hub is a hub store everything needs for create a websocket server
type Hub struct {
	MapSchedules schedule.Schedules
	Schedules    chan schedule.Schedules
	Clients      map[*websocket.Conn]bool
	Upgrader     websocket.Upgrader
}

// New create a new hub
func New() *Hub {
	return &Hub{
		MapSchedules: schedule.New(),
		Clients:      make(map[*websocket.Conn]bool),
		Schedules:    make(chan schedule.Schedules),
		Upgrader:     websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}
}

// HandleConnections is a handler for connection
func (hub *Hub) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := hub.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer ws.Close()
	hub.Clients[ws] = true

	for {
		hub.Schedules <- hub.MapSchedules
	}
}

// BoardcastSchedule boardcast MapSchedules to all clients
func (hub *Hub) BoardcastSchedule() {
	for client := range hub.Clients {
		err := client.WriteJSON(<-hub.Schedules)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(hub.Clients, client)
		}
	}
}
