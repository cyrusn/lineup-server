package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	helper "github.com/cyrusn/goHTTPHelper"
	auth "github.com/cyrusn/goJWTAuthHelper"
	jwt "github.com/dgrijalva/jwt-go"
)

var lifeTime int64 = 30

// AuthStore contains all method for handling authentication
type AuthStore interface {
	Authenticate(string, string) error
}

// AuthClaims is jwt.Claims for authentication credential
type AuthClaims struct {
	UserAlias string
	jwt.StandardClaims
}

// Update use to refresh token
func (claims *AuthClaims) Update(token *jwt.Token) {
	claims.ExpiresAt = expiresAfter(lifeTime)
	token.Claims = claims
}

// LoginHandler handle login event, will send jwt token in response
// send a POST request with JSON form
func LoginHandler(a AuthStore, s auth.Secret) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userAlias, password, err := parseJSONPostForm(r)
		errCode := http.StatusBadRequest
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}
		if err := a.Authenticate(userAlias, password); err != nil {
			helper.PrintError(w, err, errCode)
			return
		}

		claims := AuthClaims{
			UserAlias: userAlias,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expiresAfter(lifeTime),
			},
		}
		token, err := s.GenerateToken(claims)
		if err != nil {
			helper.PrintError(w, err, errCode)
			return
		}
		w.Write([]byte(token))
	}
}

// RefreshHandler refresh jwt token by a GET request
func RefreshHandler(a AuthStore, s auth.Secret) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getTokenFromHeader(r, s)
		fmt.Println(token)
		authClaims := AuthClaims{}
		newToken, err := s.UpdateToken(token, &authClaims)
		if err != nil {
			errCode := http.StatusBadRequest
			helper.PrintError(w, err, errCode)
			return
		}

		w.Write([]byte(newToken))
	}
}

func expiresAfter(min int64) int64 {
	return time.Now().Add(time.Minute * time.Duration(min)).Unix()
}

func getTokenFromHeader(r *http.Request, s auth.Secret) string {
	return r.Header.Get(s.JWTKeyName)
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
