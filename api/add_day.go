package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
)

func AddDay(w http.ResponseWriter, r *http.Request) error {
	request := service.AddDayRequest{}
	var err error
	var body []byte

	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of add day: " + err.Error())
		return err
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of add day: " + err.Error())
		return err
	}

	err = service.AddDay(&request)
	if err != nil {
		fmt.Println("An error occurred during adding day " + string(body) + ": " + err.Error())
		return err
	}

	return nil
}
