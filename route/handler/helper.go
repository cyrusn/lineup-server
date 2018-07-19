package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/cyrusn/goHTTPHelper"
	"github.com/gorilla/mux"
)

func readQueries(r *http.Request, key string) []string {
	q := r.URL.Query()
	return q[key]
}

func readQuery(r *http.Request, key string) string {
	q := r.URL.Query()
	return q.Get(key)
}

// readClassCodeAndClassNo read classcode and classno in mux.Vars
func readClassCodeAndClassNo(w http.ResponseWriter, r *http.Request) (string, int, error) {
	errCode := http.StatusBadRequest
	classCode := mux.Vars(r)["classcode"]
	classNoString := mux.Vars(r)["classno"]
	classNo, err := strconv.Atoi(classNoString)
	if err != nil {
		helper.PrintError(w, err, errCode)
		return "", 0, err
	}

	return strings.ToUpper(classCode), classNo, nil
}

func readPriority(r *http.Request) (int, error) {
	priorityString := mux.Vars(r)["priority"]
	priority, err := strconv.Atoi(priorityString)
	if err != nil {
		return 0, err
	}
	return priority, nil
}
