package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
)

func SetBreakStartDate(w http.ResponseWriter, r *http.Request) error {
	request := service.SetBreakStartDateRequest{}
	var err error
	var body []byte

	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of set start date: " + err.Error())
		return err
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of set start date: " + err.Error())
		return err
	}

	err = service.SetBreakStartDate(&request)
	if err != nil {
		fmt.Println("An error occurred during setting start date " + string(body) + ": " + err.Error())
		return err
	}

	return nil
}
