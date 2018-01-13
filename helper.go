package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	helper "github.com/cyrusn/goHTTPHelper"
	"github.com/gorilla/mux"
)

// ReadClassCodeAndClassNo read classcode and classno in mux.Vars
func ReadClassCodeAndClassNo(w http.ResponseWriter, r *http.Request) (string, int, error) {
	errCode := http.StatusBadRequest
	classCode := mux.Vars(r)["classcode"]
	classNoString := mux.Vars(r)["classno"]
	classNo, err := strconv.Atoi(classNoString)
	if err != nil {
		helper.PrintError(w, err, errCode)
		return "", 0, err
	}

	return classCode, classNo, nil
}

func readOrder(r *http.Request) (int, error) {
	orderString := mux.Vars(r)["order"]
	order, err := strconv.Atoi(orderString)
	if err != nil {
		return 0, err
	}
	return order, nil
}

func findSchedule(w http.ResponseWriter, classCode string, classNo int, print bool) *Schedule {
	for _, p := range schedules[classCode] {
		if p.ClassNo == classNo {
			return p
		}
	}
	if print {
		printUserNotFound(w, classCode, classNo)
	}
	return nil
}

func printUserNotFound(w http.ResponseWriter, classCode string, classNo int) {
	errCode := http.StatusBadRequest
	message := fmt.Sprintf("%s%d not found", classCode, classNo)
	helper.PrintError(w, errors.New(message), errCode)
}
