package api

import (
	"net/http"
)

func (a *API) GetChannels(w http.ResponseWriter, r *http.Request) (any, error) {
	return a.Service.GetChannels()
}
