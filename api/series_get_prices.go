package api

import (
	"net/http"
	"strconv"

	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
)

func (a *API) SeriesGetPrices(w http.ResponseWriter, r *http.Request) (any, error) {
	id, err := strconv.ParseInt(r.PathValue("series_id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return a.Service.SeriesGetPrices(&service.SeriesGetPricesRequest{SeriesId: id})
}
