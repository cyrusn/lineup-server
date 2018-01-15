package route

import (
	"net/http"
	"strconv"

	"github.com/cyrusn/goHTTPHelper"
	"github.com/gorilla/mux"
)

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
