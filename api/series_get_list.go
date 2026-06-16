package api

import "net/http"

func (a *API) SeriesGetList(w http.ResponseWriter, r *http.Request) (any, error) {
	return a.Service.SeriesGetList()
}
