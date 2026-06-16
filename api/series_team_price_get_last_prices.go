package api

import "net/http"

func (a *API) SeriesTeamPriceGetLastPrices(w http.ResponseWriter, r *http.Request) (any, error) {
	return a.Service.SeriesTeamPriceGetLastPrices()
}
