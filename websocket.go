package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var schedules = make(map[string][]*Schedule)
var clients = make(map[*websocket.Conn]bool)
var boardcast = make(chan map[string][]*Schedule)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer ws.Close()
	clients[ws] = true

	for {
		boardcast <- schedules
	}
}

func handleBoardcast() {
	schedule := <-boardcast
	for client := range clients {
		err := client.WriteJSON(schedule)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
