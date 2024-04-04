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

type Response struct {
	Data  any
	Error string
}

func (rb *RouteBuilder) WrapRoute(route func(w http.ResponseWriter, r *http.Request) (any, error), method string, isPrivate bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != HttpAny && r.Method != method {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		//if isPrivate {
		//	err := rb.authenticate(r.BasicAuth())
		//	if err != nil {
		//		fmt.Println("Failed auth attempt from " + r.Host)
		//		formatError(w, err, http.StatusUnauthorized)
		//		return
		//	}
		//}

		response, err := route(w, r)
		requestResponse := &Response{}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			requestResponse.Error = err.Error()
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if response != nil {
			requestResponse.Data = response
		}

		encodedResponse, errM := json.Marshal(requestResponse)
		if errM != nil {
			fmt.Println(errM)
		}

		w.Write(encodedResponse)
	}
}

type ResponseError struct {
	Error string `json:"error"`
}

func (rb *RouteBuilder) authenticate(username string, password string, ok bool) error {
	if username != rb.Username || password != rb.Password || !ok {
		return errors.New("invalid credentials")
	}

	return nil
}
