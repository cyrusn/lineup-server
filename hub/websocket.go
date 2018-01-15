package hub

import (
	"log"
	"net/http"

	"github.com/cyrusn/lineup-system/schedule"
	"github.com/gorilla/websocket"
)

// Hub is a hub store everything needs for create a websocket server
type Hub struct {
	MapSchedules     schedule.MapSchedules
	ChSchedules      chan schedule.MapSchedules
	SchedulesClients map[*websocket.Conn]bool
	Upgrader         websocket.Upgrader
}

// New create a new hub
func New() *Hub {
	return &Hub{
		MapSchedules:     schedule.New(),
		SchedulesClients: make(map[*websocket.Conn]bool),
		ChSchedules:      make(chan schedule.MapSchedules),
		Upgrader:         websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}
}

// HandleScheduleConnections is a handler for connection
func (hub *Hub) HandleScheduleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := hub.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer ws.Close()
	hub.SchedulesClients[ws] = true

	for {
		hub.ChSchedules <- hub.MapSchedules
	}
}

// BoardcastSchedule boardcast MapSchedules to all clients
func (hub *Hub) BoardcastSchedule() {
	for client := range hub.SchedulesClients {
		err := client.WriteJSON(<-hub.ChSchedules)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(hub.SchedulesClients, client)
		}
	}
}
