package api

import (
	"net/http"
)

func (a *API) GetStreams(w http.ResponseWriter, r *http.Request) (any, error) {
	return a.Service.GetStreams()
}
