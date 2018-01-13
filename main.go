package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	helper "github.com/cyrusn/goHTTPHelper"
	"github.com/gorilla/mux"
)

// Schedule ...
type Schedule struct {
	ClassNo    int       `json:"classno"`
	ArrivedAt  time.Time `json:"arrivedAt"`
	Order      int       `json:"order"`
	IsNotified bool      `json:"isNotified"`
	IsMeeting  bool      `json:"isMeeting"`
	IsComplete bool      `json:"isComplete"`
}

var (
	schedules = make(map[string][]*Schedule)
	port      string
)

func init() {
	flag.StringVar(&port, "port", ":5000", "Port value")
	flag.Parse()
}

func main() {
	r := mux.NewRouter()
	hub := newHub()
	r.HandleFunc("/ws", hub.handleConnections)

	for _, route := range hub.routes() {
		r.
			PathPrefix("/api/").
			Methods(route.Methods...).
			Path(route.Path).
			HandlerFunc(route.Handler)
	}

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../static/dist"))))

	location := "localhost" + port
	fmt.Printf("Server start on http://%s\n", location)
	http.ListenAndServe(location, helper.Logger(r))
}
