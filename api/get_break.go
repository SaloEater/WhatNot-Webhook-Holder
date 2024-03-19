package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
)

func GetBreak(w http.ResponseWriter, r *http.Request) error {
	request := service.GetBreakRequest{}
	var err error
	var body []byte

	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of get break: " + err.Error())
		return err
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of get break: " + err.Error())
		return err
	}

	response, err := service.GetBreak(&request)
	if err != nil {
		fmt.Println("An error occurred during getting break " + string(body) + ": " + err.Error())
		return err
	}

	w.Write(response)

	return nil
}
