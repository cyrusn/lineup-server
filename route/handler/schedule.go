package handler

import (
	"fmt"
	"net/http"

	"github.com/cyrusn/goHTTPHelper"
	"github.com/cyrusn/lineup-system/schedule"
)

type successMessage struct {
	Message string `json:"message"`
}

// ScheduleStore contains all method to manipulate schedule
type ScheduleStore interface {
	Insert(string, int) error
	Delete(string, int) error
	SelectByClassCode(string) ([]*schedule.Schedule, error)
	UpdatePriority(string, int, int) error
	ToggleIsNotified(string, int) error
	ToggleIsMeeting(string, int) error
	ToggleIsComplete(string, int) error
}

// GetScheduleHandler is handler to get schedules by given classcode in get request
func GetScheduleHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode := readClassCode(w, r)
		schedules, err := s.SelectByClassCode(classCode)
		errCode := http.StatusBadRequest

		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}
		helper.PrintJSON(w, schedules)
	}
}

// AddScheduleHandler is handler to add schedules by given classcode
// and classno in post request
func AddScheduleHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errCode := http.StatusBadRequest

		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		if err := s.Insert(classCode, classNo); err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		message := fmt.Sprintf("%s%d is added", classCode, classNo)
		helper.PrintJSON(w, successMessage{message})
	}
}

// RemoveScheduleHandler is handler to remove schedules by given classcode
// and classno in DELETE request
func RemoveScheduleHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errCode := http.StatusBadRequest

		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		if err := s.Delete(classCode, classNo); err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		message := fmt.Sprintf("%s%d is removed", classCode, classNo)
		helper.PrintJSON(w, successMessage{message})
	}
}

// UpdatePriorityHandler is handler to update schedules's priority by given classcode
// and classno in PUT request
func UpdatePriorityHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errCode := http.StatusBadRequest
		priority, err := readRriority(r)
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		if err := s.UpdatePriority(classCode, classNo, priority); err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		message := fmt.Sprintf("%s%d's priority updated to %d", classCode, classNo, priority)
		helper.PrintJSON(w, successMessage{message})
	}
}

// ToggleIsCompleteHandler is handler to TOGGLE schedules's IsComplete by given
// classcode and classno in PUT request
func ToggleIsCompleteHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errCode := http.StatusBadRequest
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		if err := s.ToggleIsComplete(classCode, classNo); err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		message := fmt.Sprintf("%s%d toggled completed", classCode, classNo)
		helper.PrintJSON(w, successMessage{message})
	}
}

// ToggleIsNotifiedHandler is handler to TOGGLE schedules's IsNotified by given
// classcode and classno in PUT request
func ToggleIsNotifiedHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errCode := http.StatusBadRequest
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		if err := s.ToggleIsNotified(classCode, classNo); err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		message := fmt.Sprintf("%s%d toggled isNotified", classCode, classNo)
		helper.PrintJSON(w, successMessage{message})
	}
}

// ToggleIsMeetingHandler is handler to TOGGLE schedules's IsMeeting by given
// classcode and classno in PUT request
func ToggleIsMeetingHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errCode := http.StatusBadRequest
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		if err := s.ToggleIsMeeting(classCode, classNo); err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		message := fmt.Sprintf("%s%d toggled isMeeting", classCode, classNo)
		helper.PrintJSON(w, successMessage{message})
	}
}
