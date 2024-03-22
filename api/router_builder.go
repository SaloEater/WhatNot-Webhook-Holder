package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RouteBuilder struct {
	Username string
	Password string
}

func (rb *RouteBuilder) WrapRoute(route func(w http.ResponseWriter, r *http.Request) error, method string, isPrivate bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != HttpAny && r.Method != method {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if isPrivate {
			err := rb.authenticate(r.BasicAuth())
			if err != nil {
				fmt.Println("Failed auth attempt from " + r.Host)
				formatError(w, err, http.StatusUnauthorized)
				return
			}
		}

		err := route(w, r)
		if err != nil {
			formatError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type ResponseError struct {
	Error string `json:"error"`
}

func formatError(w http.ResponseWriter, err error, code int) {
	responseError := ResponseError{Error: err.Error()}
	encoded, err := json.Marshal(responseError)
	fmt.Println(string(encoded))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(encoded)
}

func (rb *RouteBuilder) authenticate(username string, password string, ok bool) error {
	if username != rb.Username || password != rb.Password || !ok {
		return errors.New("invalid credentials")
	}

	return nil
}
