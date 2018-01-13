package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	helper "github.com/cyrusn/goHTTPHelper"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

var mapSchedule = make(map[string][]*Schedule)
var mapClient = make(map[*websocket.Conn]bool)
var boardcast = make(chan map[string][]*Schedule)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Route ...
type Route struct {
	Path    string
	Methods []string
	Handler func(http.ResponseWriter, *http.Request)
}

var routes = []Route{
	Route{
		Path:    "/schedule",
		Methods: []string{"GET"},
		Handler: getSchedule,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}",
		Methods: []string{"POST"},
		Handler: addSchedule,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}",
		Methods: []string{"DELETE"},
		Handler: removeSchedule,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}/order/{order}",
		Methods: []string{"PUT"},
		Handler: updateOrder,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}/is-complete",
		Methods: []string{"PUT"},
		Handler: toggleIsComplete,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}/is-notified",
		Methods: []string{"PUT"},
		Handler: toggleIsNotified,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}/is-meeting",
		Methods: []string{"PUT"},
		Handler: toggleIsMeeting,
	},
}

var port string

func init() {
	flag.StringVar(&port, "port", ":5000", "Port value")
	flag.Parse()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", handleConnections)

	for _, route := range routes {
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

func getSchedule(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("boardcast"))
	go handleBoardcast()
}

func addSchedule(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest

	classCode, classNo, err := ReadClassCodeAndClassNo(r)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}

	for _, p := range mapSchedule[classCode] {
		if p.ClassNo == classNo {
			message := fmt.Sprintf("%s%d already existed", classCode, classNo)
			w.Write([]byte(message))
			return
		}
	}

	mapSchedule[classCode] = append(mapSchedule[classCode], &Schedule{
		ClassNo:    classNo,
		ArrivedAt:  time.Now(),
		Order:      0,
		IsNotified: false,
		IsMeeting:  false,
		IsComplete: false,
	})

	message := fmt.Sprintf("%s%d added", classCode, classNo)
	w.Write([]byte(message))

	go handleBoardcast()
}

func removeSchedule(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest
	classCode, classNo, err := ReadClassCodeAndClassNo(r)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}

	var newList = []*Schedule{}

	for _, p := range mapSchedule[classCode] {
		if p.ClassNo == classNo {
			continue
		} else {
			newList = append(newList, p)
		}
	}
	if len(mapSchedule[classCode]) == len(newList) {
		message := fmt.Sprintf("%s%d not found", classCode, classNo)
		w.Write([]byte(message))
		return
	}

	mapSchedule[classCode] = newList
	message := fmt.Sprintf("%s%d removed", classCode, classNo)
	w.Write([]byte(message))

	go handleBoardcast()
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest
	orderString := mux.Vars(r)["order"]
	order, err := strconv.Atoi(orderString)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}

	classCode, classNo, err := ReadClassCodeAndClassNo(r)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}

	for _, p := range mapSchedule[classCode] {
		if p.ClassNo == classNo {
			p.Order = order
			message := fmt.Sprintf("%s%d updated order to %d", classCode, classNo, order)
			w.Write([]byte(message))
			go handleBoardcast()
			return
		}
	}

	message := fmt.Sprintf("%s%d not found", classCode, classNo)
	w.Write([]byte(message))

}

func toggleIsComplete(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest
	classCode, classNo, err := ReadClassCodeAndClassNo(r)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}

	for _, p := range mapSchedule[classCode] {
		if p.ClassNo == classNo {
			p.IsComplete = !p.IsComplete

			go handleBoardcast()
			message := fmt.Sprintf("%s%d toggled completed to %v", classCode, classNo, p.IsComplete)
			w.Write([]byte(message))
			return
		}
	}

	message := fmt.Sprintf("%s%d not found", classCode, classNo)
	w.Write([]byte(message))

}

func toggleIsNotified(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest

	classCode, classNo, err := ReadClassCodeAndClassNo(r)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}

	for _, p := range mapSchedule[classCode] {
		if p.ClassNo == classNo {
			p.IsNotified = !p.IsNotified

			go handleBoardcast()
			message := fmt.Sprintf("%s%d toggled isNotified to %v", classCode, classNo, p.IsNotified)
			w.Write([]byte(message))
			return
		}
	}

	message := fmt.Sprintf("%s%d not found", classCode, classNo)
	w.Write([]byte(message))

}

func toggleIsMeeting(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest

	classCode, classNo, err := ReadClassCodeAndClassNo(r)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}

	for _, p := range mapSchedule[classCode] {
		if p.ClassNo == classNo {
			p.IsMeeting = !p.IsMeeting

			message := fmt.Sprintf("%s%d toggled isMeeting to %v", classCode, classNo, p.IsNotified)
			w.Write([]byte(message))
			go handleBoardcast()
			return
		}
	}

	message := fmt.Sprintf("%s%d not found", classCode, classNo)
	w.Write([]byte(message))

}

// ReadClassCodeAndClassNo read classcode and classno in mux.Vars
func ReadClassCodeAndClassNo(r *http.Request) (string, int, error) {
	classCode := mux.Vars(r)["classcode"]
	classNoString := mux.Vars(r)["classno"]
	classNo, err := strconv.Atoi(classNoString)
	if err != nil {
		return "", 0, err
	}

	return classCode, classNo, nil
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer ws.Close()
	mapClient[ws] = true

	for {
		boardcast <- mapSchedule
	}
}

func handleBoardcast() {
	schedule := <-boardcast
	for client := range mapClient {
		err := client.WriteJSON(schedule)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(mapClient, client)
		}
	}
}
