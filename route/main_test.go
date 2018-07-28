package route_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cyrusn/goTestHelper"
	jwt "github.com/dgrijalva/jwt-go"

	auth_helper "github.com/cyrusn/goJWTAuthHelper"

	"github.com/cyrusn/lineup-system/model/auth"
	"github.com/cyrusn/lineup-system/model/schedule"
	"github.com/cyrusn/lineup-system/route"

	"github.com/gorilla/mux"
)

var r = mux.NewRouter()
var secret auth_helper.Secret

type mockAuthDB struct {
	credentials []*auth.Credential
}

type myClaims struct {
	jwt.StandardClaims
}

func (db *mockAuthDB) Authenticate(username, password string) (string, error) {
	if username == "user1" && password == "password1" {
		return secret.GenerateToken(myClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * time.Duration(30)).Unix(),
			},
		})
	}
	return "", errors.New("invalid login")
}

func (db *mockAuthDB) Refresh(jwtToken string) (string, error) {
	return secret.GenerateToken(myClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(30)).Unix(),
		},
	})
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

func (db *mockScheduleDB) SelectedBy(q *schedule.Query) ([]*schedule.Schedule, error) {
	var result []*schedule.Schedule

	for _, s := range db.schedules {
		for _, classcode := range q.Classcodes {
			if s.ClassCode == classcode {
				result = append(result, s)
			}
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
			&auth.Credential{"user1", "password1", "teacher"},
			&auth.Credential{"user2", "password2", "student"},
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
}

func init() {
	secret = auth_helper.New("myClaim", "kid", "myRole", []byte("secret"))
	for _, route := range route.Routes(&db) {
		handler := http.HandlerFunc(route.Handler)

		if route.Auth {
			handler = secret.Authenticate(handler).(http.HandlerFunc)
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
