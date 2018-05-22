package route_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cyrusn/goTestHelper"

	"github.com/cyrusn/lineup-system/model"
	"github.com/cyrusn/lineup-system/route"

	"github.com/gorilla/mux"
)

var r = mux.NewRouter()

type mockScheduleDB struct {
	schedules []*model.Schedule
}

func (db *mockScheduleDB) Insert(classCode string, classNo int) error {
	var newSchedule = &model.Schedule{
		classCode, classNo, time.Now(), 0, false, false, false,
	}
	db.schedules = append(db.schedules, newSchedule)
	return nil
}

func (db *mockScheduleDB) Delete(classCode string, classNo int) error {
	var index int
	for i, s := range db.schedules {
		if s.ClassCode == classCode && s.ClassNo == classNo {
			index = i
		}
	}

	db.schedules = append(db.schedules[:index], db.schedules[index+1:]...)
	return nil
}

func (db *mockScheduleDB) SelectByClassCode(classCode string) ([]*model.Schedule, error) {
	var result []*model.Schedule

	for _, s := range db.schedules {
		if s.ClassCode == classCode {
			result = append(result, s)
		}
	}
	return result, nil
}

func (db *mockScheduleDB) UpdatePriority(classCode string, classNo, priority int) error {
	for _, s := range db.schedules {
		if s.ClassCode == classCode && s.ClassNo == classNo {
			s.Priority = priority
		}
	}
	return nil
}

func (db *mockScheduleDB) ToggleIsNotified(classCode string, classNo int) error {
	for _, s := range db.schedules {
		if s.ClassCode == classCode && s.ClassNo == classNo {
			s.IsNotified = !s.IsNotified
		}
	}
	return nil
}

func (db *mockScheduleDB) ToggleIsMeeting(classCode string, classNo int) error {
	for _, s := range db.schedules {
		if s.ClassCode == classCode && s.ClassNo == classNo {
			s.IsMeeting = !s.IsMeeting
		}
	}
	return nil
}

func (db *mockScheduleDB) ToggleIsComplete(classCode string, classNo int) error {
	for _, s := range db.schedules {
		if s.ClassCode == classCode && s.ClassNo == classNo {
			s.IsComplete = !s.IsComplete
		}
	}
	return nil
}

var db = mockScheduleDB{
	[]*model.Schedule{
		&model.Schedule{"3A", 12, time.Now(), 0, false, false, false},
		&model.Schedule{"3C", 8, time.Now(), 0, false, false, false},
		&model.Schedule{"3B", 2, time.Now(), 0, false, false, false},
		&model.Schedule{"3C", 5, time.Now(), 0, false, false, false},
		&model.Schedule{"3D", 1, time.Now(), 0, false, false, false},
		&model.Schedule{"3B", 4, time.Now(), 0, false, false, false},
		&model.Schedule{"3A", 9, time.Now(), 0, false, false, false},
		&model.Schedule{"3D", 10, time.Now(), 0, false, false, false},
		&model.Schedule{"3C", 14, time.Now(), 0, false, false, false},
		&model.Schedule{"3C", 6, time.Now(), 0, false, false, false},
		&model.Schedule{"3D", 2, time.Now(), 0, false, false, false},
	},
}

func init() {
	for _, route := range route.Routes(&db) {
		handler := http.HandlerFunc(route.Handler)
		r.Handle(route.Path, handler).Methods(route.Methods...)
	}
}

// Helper functions
func parseBody(w *httptest.ResponseRecorder, t *testing.T) []byte {
	resp := w.Result()
	if resp.StatusCode >= 400 {
		assert.OK(t, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	assert.OK(t, err)
	return body
}
