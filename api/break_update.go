package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
)

func (a *API) UpdateBreak(w http.ResponseWriter, r *http.Request) (any, error) {
	request := service.UpdateBreakRequest{}
	var err error
	var body []byte

	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of update break: " + err.Error())
		return nil, err
	}

	b := string(body)
	fmt.Println(b)
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of update break: " + err.Error())
		return nil, err
	}

	return a.Service.UpdateBreak(&request)
}
