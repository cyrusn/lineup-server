package route

import (
	"fmt"
	"net/http"

	"github.com/cyrusn/goHTTPHelper"
	"github.com/cyrusn/lineup-system/model"
)

type successMessage struct {
	Message string `json:"message"`
}

// ScheduleStore contains all method to manipulate schedule
type ScheduleStore interface {
	Insert(string, int) error
	Delete(string, int) error
	SelectByClassCode(string) ([]*model.Schedule, error)
	UpdatePriority(string, int, int) error
	ToggleIsNotified(string, int) error
	ToggleIsMeeting(string, int) error
	ToggleIsComplete(string, int) error
}

func getScheduleHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode := readClassCode(w, r)
		schedules, err := s.SelectByClassCode(classCode)
		errCode := http.StatusBadRequest

		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}
		helper.PrintJSON(w, schedules, errCode)
	}
}

func addScheduleHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
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
		helper.PrintJSON(w, successMessage{message}, errCode)
	}
}

func removeScheduleHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
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
		helper.PrintJSON(w, successMessage{message}, errCode)
	}
}

func updatePriorityHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
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
		helper.PrintJSON(w, successMessage{message}, errCode)
	}
}

func toggleIsCompleteHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
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
		helper.PrintJSON(w, successMessage{message}, errCode)
	}
}

func toggleIsNotifiedHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
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
		helper.PrintJSON(w, successMessage{message}, errCode)
	}
}

func toggleIsMeetingHandler(s ScheduleStore) func(http.ResponseWriter, *http.Request) {
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
		helper.PrintJSON(w, successMessage{message}, errCode)
	}
}
