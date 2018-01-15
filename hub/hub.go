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
	SchedulesClients map[*websocket.Conn]bool
}

// New create a new hub
func New() *Hub {
	return &Hub{
		MapSchedules:     schedule.New(),
		SchedulesClients: make(map[*websocket.Conn]bool),
	}
}

// ServeWS is a handler for connection
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	Upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer ws.Close()
	h.SchedulesClients[ws] = true
	ReadLoop(ws)
}

// ReadLoop read the connection to process close, ping and pong messages
func ReadLoop(ws *websocket.Conn) {
	for {
		if _, _, err := ws.NextReader(); err != nil {
			ws.Close()
			break
		}
	}
}

// BoardcastSchedule boardcast MapSchedules to all clients
func (h *Hub) BoardcastSchedule() {
	for client := range h.SchedulesClients {
		if err := client.WriteJSON(h.MapSchedules); err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(h.SchedulesClients, client)
		}
	}
}
