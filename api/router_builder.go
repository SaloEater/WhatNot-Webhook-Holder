package api

import (
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
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		err := route(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (rb *RouteBuilder) authenticate(username string, password string, ok bool) error {
	if username != rb.Username || password != rb.Password || !ok {
		return errors.New("invalid credentials")
	}

	return nil
}
