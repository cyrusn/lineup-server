package route_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	auth_helper "github.com/cyrusn/goJWTAuthHelper"
	"github.com/cyrusn/goTestHelper"

	auth "github.com/cyrusn/lineup-system/auth"
	"github.com/cyrusn/lineup-system/route"
	"github.com/cyrusn/lineup-system/schedule"

	"github.com/gorilla/mux"
)

var r = mux.NewRouter()

type mockAuthDB struct {
	credentials []*auth.Credential
}

func (db *mockAuthDB) Authenticate(username, password string) error {
	return nil
}
func (db *mockAuthDB) Validate(token string) error {
	return nil
}

type mockScheduleDB struct {
	schedules []*schedule.Schedule
}

func (db *mockScheduleDB) Insert(classCode string, classNo int) error {
	var newSchedule = &schedule.Schedule{
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

func (db *mockScheduleDB) SelectByClassCode(classCode string) ([]*schedule.Schedule, error) {
	var result []*schedule.Schedule

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

var db = route.Store{
	&mockAuthDB{
		[]*auth.Credential{
			&auth.Credential{"user1", "password1"},
			&auth.Credential{"user2", "password2"},
		},
	},
	&mockScheduleDB{
		[]*schedule.Schedule{
			&schedule.Schedule{"3A", 12, time.Now(), 0, false, false, false},
			&schedule.Schedule{"3C", 8, time.Now(), 0, false, false, false},
			&schedule.Schedule{"3B", 2, time.Now(), 0, false, false, false},
			&schedule.Schedule{"3C", 5, time.Now(), 0, false, false, false},
			&schedule.Schedule{"3D", 1, time.Now(), 0, false, false, false},
			&schedule.Schedule{"3B", 4, time.Now(), 0, false, false, false},
			&schedule.Schedule{"3A", 9, time.Now(), 0, false, false, false},
			&schedule.Schedule{"3D", 10, time.Now(), 0, false, false, false},
			&schedule.Schedule{"3C", 14, time.Now(), 0, false, false, false},
			&schedule.Schedule{"3C", 6, time.Now(), 0, false, false, false},
			&schedule.Schedule{"3D", 2, time.Now(), 0, false, false, false},
		},
	},
	auth_helper.New("myClaim", "kid", "myRole", []byte("secret")),
}

func init() {
	for _, route := range route.Routes(&db) {
		handler := http.HandlerFunc(route.Handler)

		if route.Auth {
			handler = db.Secret.Authenticate(handler).(http.HandlerFunc)
		}

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
