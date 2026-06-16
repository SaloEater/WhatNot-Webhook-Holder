package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
)

func (a *API) SeriesCreate(w http.ResponseWriter, r *http.Request) (any, error) {
	request := service.SeriesCreateRequest{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of series create: " + err.Error())
		return nil, err
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of series create: " + err.Error())
		return nil, err
	}
	return a.Service.SeriesCreate(&request)
}
