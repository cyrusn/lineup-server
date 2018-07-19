package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	helper "github.com/cyrusn/goHTTPHelper"
	"github.com/gorilla/mux"
)

// AuthStore contains all method for handling authentication
type AuthStore interface {
	Authenticate(string, string) (string, error)
	Refresh(string) (string, error)
}

// LoginHandler handle login event, will send jwt token in response
// send a POST request with JSON form
func LoginHandler(a AuthStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userAlias, password, err := parseJSONPostForm(r)
		errCode := http.StatusBadRequest
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}
		token, err := a.Authenticate(userAlias, password)
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		w.Write([]byte(token))
	}
}

// RefreshHandler refresh jwt token by a GET request,
func RefreshHandler(a AuthStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := mux.Vars(r)["key"]
		newToken, err := a.Refresh(token)

		if err != nil {
			errCode := http.StatusBadRequest
			helper.PrintError(w, err, errCode)
			return
		}

		w.Write([]byte(newToken))
	}
}

func parseJSONPostForm(r *http.Request) (string, string, error) {
	type loginForm struct {
		UserAlias string
		Password  string
	}
	form := new(loginForm)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", "", err
	}

	if err = json.Unmarshal(body, form); err != nil {
		return "", "", err
	}

	return form.UserAlias, form.Password, nil
}
