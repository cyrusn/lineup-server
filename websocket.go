package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Hub ...
type Hub struct {
	Clients   map[*websocket.Conn]bool
	Schedules chan map[string][]*Schedule
}

func newHub() *Hub {
	return &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Schedules: make(chan map[string][]*Schedule),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (hub *Hub) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer ws.Close()
	hub.Clients[ws] = true

	for {
		hub.Schedules <- schedules
	}
}

func (hub *Hub) boardcast() {
	schedule := <-hub.Schedules
	for client := range hub.Clients {
		err := client.WriteJSON(schedule)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(hub.Clients, client)
		}
	}
}
