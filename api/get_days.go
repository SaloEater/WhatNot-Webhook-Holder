package api

import (
	"net/http"
)

func (a *API) GetDays(w http.ResponseWriter, r *http.Request) (any, error) {
	return a.Service.GetDays()
}
