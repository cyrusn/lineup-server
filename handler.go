package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	helper "github.com/cyrusn/goHTTPHelper"
)

func (hub *Hub) getScheduleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("boardcast"))
	go hub.boardcast()
}

func (hub *Hub) addScheduleHandler(w http.ResponseWriter, r *http.Request) {
	classCode, classNo, err := ReadClassCodeAndClassNo(w, r)
	if err != nil {
		return
	}

	p := findSchedule(w, classCode, classNo, false)
	if p != nil {
		errCode := http.StatusBadRequest
		message := fmt.Sprintf("%s%d is already exist", classCode, classNo)
		helper.PrintError(w, errors.New(message), errCode)
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

	go hub.boardcast()
}

func (hub *Hub) removeScheduleHandler(w http.ResponseWriter, r *http.Request) {
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
	go hub.boardcast()

	message := fmt.Sprintf("%s%d removed", classCode, classNo)
	w.Write([]byte(message))
}

func (hub *Hub) updateOrderHandler(w http.ResponseWriter, r *http.Request) {
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
	go hub.boardcast()

	message := fmt.Sprintf("%s%d updated order to %d", classCode, classNo, order)
	w.Write([]byte(message))
}

func (hub *Hub) toggleIsCompleteHandler(w http.ResponseWriter, r *http.Request) {
	classCode, classNo, err := ReadClassCodeAndClassNo(w, r)
	if err != nil {
		return
	}

	p := findSchedule(w, classCode, classNo, true)
	if p == nil {
		return
	}

	p.IsComplete = !p.IsComplete
	go hub.boardcast()

	message := fmt.Sprintf("%s%d toggled completed to %v", classCode, classNo, p.IsComplete)
	w.Write([]byte(message))
}

func (hub *Hub) toggleIsNotifiedHandler(w http.ResponseWriter, r *http.Request) {
	classCode, classNo, err := ReadClassCodeAndClassNo(w, r)
	if err != nil {
		return
	}

	p := findSchedule(w, classCode, classNo, true)
	if p == nil {
		return
	}

	p.IsNotified = !p.IsNotified
	go hub.boardcast()

	message := fmt.Sprintf("%s%d toggled isNotified to %v", classCode, classNo, p.IsNotified)
	w.Write([]byte(message))
}

func (hub *Hub) toggleIsMeetingHandler(w http.ResponseWriter, r *http.Request) {
	classCode, classNo, err := ReadClassCodeAndClassNo(w, r)
	if err != nil {
		return
	}

	p := findSchedule(w, classCode, classNo, true)
	if p == nil {
		return
	}

	p.IsMeeting = !p.IsMeeting
	go hub.boardcast()

	message := fmt.Sprintf("%s%d toggled isMeeting to %v", classCode, classNo, p.IsNotified)
	w.Write([]byte(message))
}
