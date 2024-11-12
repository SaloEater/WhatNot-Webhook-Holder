package api

import (
	"net/http"
)

func (a *API) LocationGetList(w http.ResponseWriter, r *http.Request) (any, error) {
	return a.Service.LocationGetList()
}
