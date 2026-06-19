package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
)

func (a *API) GetWidgetSeriesPick2(w http.ResponseWriter, r *http.Request) (any, error) {
	request := service.GetWidgetSeriesPick2Request{}
	var err error
	var body []byte

	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of get widget series pick2: " + err.Error())
		return nil, err
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of get widget series pick2: " + err.Error())
		return nil, err
	}

	return a.Service.GetWidgetSeriesPick2(&request)
}
