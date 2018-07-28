package route_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

var jwtToken string

// TestMain run test
func TestMain(t *testing.T) {
	t.Run("login teacher1", login)
	t.Run("insert schedule 3A20", addSchedule("3A", 20))
	t.Run("delete schedule 3A12", deleteSchedule("3A", 12))
	t.Run("update schedule 3A20 to 2", updateSchedulePriority("3A", 20, 2))

	for _, key := range []string{"is-complete", "is-notified", "is-meeting"} {
		classsCode := "3A"
		classNo := 20
		name := fmt.Sprintf("toggle %s schedule %s%d", key, classsCode, classNo)
		t.Run(name, toggler(classsCode, classNo, key))
	}

	t.Run("get 3A user", selectScheduleByClass("3A"))
}

var login = func(t *testing.T) {
	w := httptest.NewRecorder()
	route := fmt.Sprintf("/auth/login")

	formString := fmt.Sprintf(`{"userAlias":"%s", "password":"%s"}`, "user1", "password1")
	postForm := strings.NewReader(formString)

	req := httptest.NewRequest("POST", route, postForm)
	r.ServeHTTP(w, req)
	body := parseBody(w, t)
	jwtToken = string(body)
}

var selectScheduleByClass = func(classCode string) func(*testing.T) {
	return func(t *testing.T) {
		w := httptest.NewRecorder()
		route := fmt.Sprintf("/schedules?classcode=%s", classCode)
		req := httptest.NewRequest("GET", route, nil)
		req.Header.Set("kid", jwtToken)
		r.ServeHTTP(w, req)

		body := parseBody(w, t)
		fmt.Println(string(body))
	}
}

var addSchedule = func(classcode string, classno int) func(*testing.T) {
	return func(t *testing.T) {
		w := httptest.NewRecorder()
		route := fmt.Sprintf("/schedule/%s/%d", classcode, classno)
		req := httptest.NewRequest("POST", route, nil)
		req.Header.Set("kid", jwtToken)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		body := parseBody(w, t)
		fmt.Println(string(body))
	}
}

var deleteSchedule = func(classcode string, classno int) func(*testing.T) {
	return func(t *testing.T) {
		w := httptest.NewRecorder()
		route := fmt.Sprintf("/schedule/%s/%d", classcode, classno)
		req := httptest.NewRequest("DELETE", route, nil)
		req.Header.Set("kid", jwtToken)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		body := parseBody(w, t)
		fmt.Println(string(body))
	}
}

var updateSchedulePriority = func(classcode string, classno, priority int) func(*testing.T) {
	return func(t *testing.T) {
		w := httptest.NewRecorder()
		route := fmt.Sprintf("/schedule/%s/%d/priority/%d", classcode, classno, 2)
		req := httptest.NewRequest("PUT", route, nil)
		req.Header.Set("kid", jwtToken)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		body := parseBody(w, t)
		fmt.Println(string(body))
	}
}

var toggler = func(classcode string, classno int, key string) func(*testing.T) {
	return func(t *testing.T) {
		w := httptest.NewRecorder()
		route := fmt.Sprintf("/schedule/%s/%d/%s", classcode, classno, key)
		req := httptest.NewRequest("PUT", route, nil)
		req.Header.Set("kid", jwtToken)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		body := parseBody(w, t)
		fmt.Println(string(body))
	}
}
