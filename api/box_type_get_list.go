package api

import (
	"net/http"
)

func (a *API) BoxTypeGetList(w http.ResponseWriter, r *http.Request) (any, error) {
	return a.Service.BoxTypeGetList()
}
