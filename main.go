package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	helper "github.com/cyrusn/goHTTPHelper"
	"github.com/gorilla/mux"
)

// Parents ...
type Parents struct {
	ClassNo   int
	Arrival   time.Time
	Moving    bool
	Completed bool
	Priority  int
}

var mapParents = make(map[string][]*Parents)

// Route ...
type Route struct {
	Path    string
	Methods []string
	Handler func(http.ResponseWriter, *http.Request)
}

var routes = []Route{
	Route{
		Path:    "/arrival",
		Methods: []string{"GET"},
		Handler: getArrival,
	},
	Route{
		Path:    "/arrival/{classcode}/{classno}",
		Methods: []string{"POST"},
		Handler: newArrival,
	},
	Route{
		Path:    "/arrival/{classcode}/{classno}",
		Methods: []string{"DELETE"},
		Handler: removeArrival,
	},
	Route{
		Path:    "/arrival/{classcode}/{classno}/priority/{priority}",
		Methods: []string{"PUT"},
		Handler: updateArrivalPriority,
	},
	Route{
		Path:    "/arrival/{classcode}/{classno}/complete",
		Methods: []string{"PUT"},
		Handler: toggleComplete,
	},
	Route{
		Path:    "/arrival/{classcode}/{classno}/moving",
		Methods: []string{"PUT"},
		Handler: toggleMoving,
	},
}

var port string

func init() {
	flag.StringVar(&port, "port", ":5000", "Port value")
	flag.Parse()
}
func main() {
	r := mux.NewRouter()

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

func getArrival(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest
	helper.PrintJSON(w, mapParents, errCode)
}

func newArrival(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest

	classCode := mux.Vars(r)["classcode"]
	classNoString := mux.Vars(r)["classno"]
	classNo, err := strconv.Atoi(classNoString)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}

	for _, p := range mapParents[classCode] {
		if p.ClassNo == classNo {
			message := fmt.Sprintf("%s%d already existed", classCode, classNo)
			w.Write([]byte(message))
			return
		}
	}

	mapParents[classCode] = append(mapParents[classCode], &Parents{
		ClassNo:   classNo,
		Arrival:   time.Now(),
		Priority:  0,
		Moving:    false,
		Completed: false,
	})
	message := fmt.Sprintf("%s%d added", classCode, classNo)
	w.Write([]byte(message))
}

func removeArrival(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest

	classCode := mux.Vars(r)["classcode"]
	classNoString := mux.Vars(r)["classno"]
	classNo, err := strconv.Atoi(classNoString)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}

	var newParentList = []*Parents{}

	for _, p := range mapParents[classCode] {
		if p.ClassNo == classNo {
			continue
		} else {
			newParentList = append(newParentList, p)
		}
	}
	if len(mapParents[classCode]) == len(newParentList) {
		message := fmt.Sprintf("%s%d not found", classCode, classNo)
		w.Write([]byte(message))
		return
	}

	mapParents[classCode] = newParentList
	message := fmt.Sprintf("%s%d removed", classCode, classNo)
	w.Write([]byte(message))
}

func updateArrivalPriority(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest

	classCode := mux.Vars(r)["classcode"]
	classNoString := mux.Vars(r)["classno"]
	priorityString := mux.Vars(r)["priority"]

	classNo, err := strconv.Atoi(classNoString)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}
	priority, err := strconv.Atoi(priorityString)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}
	for _, p := range mapParents[classCode] {
		if p.ClassNo == classNo {
			p.Priority = priority
			message := fmt.Sprintf("%s%d updated Priority to %d", classCode, classNo, priority)
			w.Write([]byte(message))
			return
		}
	}

	message := fmt.Sprintf("%s%d not found", classCode, classNo)
	w.Write([]byte(message))
}

func toggleComplete(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest
	classCode := mux.Vars(r)["classcode"]
	classNoString := mux.Vars(r)["classno"]
	classNo, err := strconv.Atoi(classNoString)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}

	for _, p := range mapParents[classCode] {
		if p.ClassNo == classNo {
			p.Completed = !p.Completed

			message := fmt.Sprintf("%s%d toggled completed to %v", classCode, classNo, p.Completed)
			w.Write([]byte(message))
			return
		}
	}

	message := fmt.Sprintf("%s%d not found", classCode, classNo)
	w.Write([]byte(message))
}
func toggleMoving(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusBadRequest
	classCode := mux.Vars(r)["classcode"]
	classNoString := mux.Vars(r)["classno"]
	classNo, err := strconv.Atoi(classNoString)
	if err != nil {
		helper.PrintError(w, err, errCode)
	}

	for _, p := range mapParents[classCode] {
		if p.ClassNo == classNo {
			p.Moving = !p.Moving

			message := fmt.Sprintf("%s%d toggled Moving to %v", classCode, classNo, p.Moving)
			w.Write([]byte(message))
			return
		}
	}

	message := fmt.Sprintf("%s%d not found", classCode, classNo)
	w.Write([]byte(message))
}
