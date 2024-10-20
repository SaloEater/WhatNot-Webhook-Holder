package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
)

func (a *API) EventStreamEnded(w http.ResponseWriter, r *http.Request) (any, error) {
	request := service.EventStreamEndedRequest{}
	var err error
	var body []byte

	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of EventStreamEnded: " + err.Error())
		return nil, err
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of EventStreamEnded: " + err.Error())
		return nil, err
	}

	return nil, a.Service.EventStreamEnded(&request)
}
