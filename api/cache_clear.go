package api

import (
	"net/http"
)

func (a *API) CacheClear(w http.ResponseWriter, r *http.Request) (any, error) {
	return nil, a.Service.CacheClear()
}
