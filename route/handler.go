package route

import (
	"fmt"
	"net/http"

	helper "github.com/cyrusn/goHTTPHelper"
	"github.com/cyrusn/lineup-system/hub"
)

func getScheduleHandler(hub *hub.Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("boardcast"))
		go hub.BoardcastSchedule()
	}
}

func addScheduleHandler(hub *hub.Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			return
		}

		if err := hub.MapSchedules.AppendSchedule(classCode, classNo); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}

		message := fmt.Sprintf("%s%d added", classCode, classNo)
		w.Write([]byte(message))

		go hub.BoardcastSchedule()
	}
}

func removeScheduleHandler(hub *hub.Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			return
		}

		if err := hub.MapSchedules.RemoveSchedule(classCode, classNo); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}

		go hub.BoardcastSchedule()

		message := fmt.Sprintf("%s%d removed", classCode, classNo)
		w.Write([]byte(message))
	}
}

func updateOrderHandler(hub *hub.Hub) func(http.ResponseWriter, *http.Request) {
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

		if err := hub.MapSchedules.UpdateOrder(classCode, classNo, order); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}
		go hub.BoardcastSchedule()

		message := fmt.Sprintf("%s%d updated order to %d", classCode, classNo, order)
		w.Write([]byte(message))
	}
}

func toggleIsCompleteHandler(hub *hub.Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			return
		}

		if err := hub.MapSchedules.ToggleIsComplete(classCode, classNo); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}

		go hub.BoardcastSchedule()

		message := fmt.Sprintf("%s%d toggled completed", classCode, classNo)
		w.Write([]byte(message))
	}
}

func toggleIsNotifiedHandler(hub *hub.Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			return
		}

		if err := hub.MapSchedules.ToggleIsNotified(classCode, classNo); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}
		go hub.BoardcastSchedule()

		message := fmt.Sprintf("%s%d toggled isNotified", classCode, classNo)
		w.Write([]byte(message))
	}
}

func toggleIsMeetingHandler(hub *hub.Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classCode, classNo, err := readClassCodeAndClassNo(w, r)
		if err != nil {
			return
		}

		if err := hub.MapSchedules.ToggleIsIsMeeting(classCode, classNo); err != nil {
			errCode := http.StatusNotFound
			helper.PrintError(w, err, errCode)
			return
		}
		go hub.BoardcastSchedule()

		message := fmt.Sprintf("%s%d toggled isMeeting", classCode, classNo)
		w.Write([]byte(message))
	}
}
