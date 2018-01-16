package route

import (
	"fmt"
	"net/http"

	"github.com/cyrusn/goHTTPHelper"
	"github.com/cyrusn/lineup-system/hub"
	"github.com/cyrusn/lineup-system/schedule"
)

func getScheduleHandler(hub *hub.Hub, s schedule.MapSchedules) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("boardcast"))
		hub.ChanMapScheudle <- s
	}
}

func addScheduleHandler(hub *hub.Hub, s schedule.MapSchedules) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			return
		}

		if err := s.AppendSchedule(classCode, classNo); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}

		message := fmt.Sprintf("%s%d added", classCode, classNo)
		w.Write([]byte(message))

		hub.ChanMapScheudle <- s
	}
}

func removeScheduleHandler(hub *hub.Hub, s schedule.MapSchedules) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			return
		}

		if err := s.RemoveSchedule(classCode, classNo); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}

		hub.ChanMapScheudle <- s

		message := fmt.Sprintf("%s%d removed", classCode, classNo)
		w.Write([]byte(message))
	}
}

func updateOrderHandler(hub *hub.Hub, s schedule.MapSchedules) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errCode := http.StatusBadRequest
		order, err := readOrder(r)
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			return
		}

		if err := s.UpdateOrder(classCode, classNo, order); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}
		hub.ChanMapScheudle <- s

		message := fmt.Sprintf("%s%d updated order to %d", classCode, classNo, order)
		w.Write([]byte(message))
	}
}

func toggleIsCompleteHandler(hub *hub.Hub, s schedule.MapSchedules) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			return
		}

		if err := s.ToggleIsComplete(classCode, classNo); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}

		hub.ChanMapScheudle <- s

		message := fmt.Sprintf("%s%d toggled completed", classCode, classNo)
		w.Write([]byte(message))
	}
}

func toggleIsNotifiedHandler(hub *hub.Hub, s schedule.MapSchedules) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			return
		}

		if err := s.ToggleIsNotified(classCode, classNo); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}
		hub.ChanMapScheudle <- s

		message := fmt.Sprintf("%s%d toggled isNotified", classCode, classNo)
		w.Write([]byte(message))
	}
}

func toggleIsMeetingHandler(hub *hub.Hub, s schedule.MapSchedules) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			return
		}

		if err := s.ToggleIsMeeting(classCode, classNo); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}
		hub.ChanMapScheudle <- s

		message := fmt.Sprintf("%s%d toggled isMeeting", classCode, classNo)
		w.Write([]byte(message))
	}
}
