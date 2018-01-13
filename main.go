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

func getScheduleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("boardcast"))
	go handleBoardcast()
}

func addScheduleHandler(w http.ResponseWriter, r *http.Request) {
	classCode, classNo, err := ReadClassCodeAndClassNo(w, r)
	if err != nil {
		return
	}

	p := findSchedule(w, classCode, classNo, false)
	if p != nil {
		return
	}

	schedules[classCode] = append(schedules[classCode], &Schedule{
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

func removeScheduleHandler(w http.ResponseWriter, r *http.Request) {
	classCode, classNo, err := ReadClassCodeAndClassNo(w, r)
	if err != nil {
		return
	}

	var newList = []*Schedule{}

	for _, p := range schedules[classCode] {
		if p.ClassNo == classNo {
			continue
		} else {
			newList = append(newList, p)
		}
	}
	if len(schedules[classCode]) == len(newList) {
		printUserNotFound(w, classCode, classNo)
		return
	}

	schedules[classCode] = newList
	message := fmt.Sprintf("%s%d removed", classCode, classNo)
	w.Write([]byte(message))

	go handleBoardcast()
}

func updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest
	order, err := readOrder(r)
	if err != nil {
		helper.PrintError(w, err, errCode)
		return
	}

	classCode, classNo, err := ReadClassCodeAndClassNo(w, r)
	if err != nil {
		return
	}

	p := findSchedule(w, classCode, classNo, true)
	if p == nil {
		return
	}

	p.Order = order
	go handleBoardcast()
	message := fmt.Sprintf("%s%d updated order to %d", classCode, classNo, order)
	w.Write([]byte(message))
}

func toggleIsCompleteHandler(w http.ResponseWriter, r *http.Request) {
	classCode, classNo, err := ReadClassCodeAndClassNo(w, r)
	if err != nil {
		return
	}

	p := findSchedule(w, classCode, classNo, true)
	if p == nil {
		return
	}

	p.IsComplete = !p.IsComplete

	go handleBoardcast()
	message := fmt.Sprintf("%s%d toggled completed to %v", classCode, classNo, p.IsComplete)
	w.Write([]byte(message))
}

func toggleIsNotifiedHandler(w http.ResponseWriter, r *http.Request) {
	classCode, classNo, err := ReadClassCodeAndClassNo(w, r)
	if err != nil {
		return
	}

	p := findSchedule(w, classCode, classNo, true)
	if p == nil {
		return
	}

	p.IsNotified = !p.IsNotified
	go handleBoardcast()
	message := fmt.Sprintf("%s%d toggled isNotified to %v", classCode, classNo, p.IsNotified)
	w.Write([]byte(message))
}

func toggleIsMeetingHandler(w http.ResponseWriter, r *http.Request) {
	classCode, classNo, err := ReadClassCodeAndClassNo(w, r)
	if err != nil {
		return
	}

	p := findSchedule(w, classCode, classNo, true)
	if p == nil {
		return
	}

	p.IsMeeting = !p.IsMeeting
	message := fmt.Sprintf("%s%d toggled isMeeting to %v", classCode, classNo, p.IsNotified)
	w.Write([]byte(message))
	go handleBoardcast()
}
