package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
)

func ChangeOutcome(w http.ResponseWriter, r *http.Request) error {
	request := service.AddOutcomeRequest{}
	var err error
	var body []byte

	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of add outcome: " + err.Error())
		return err
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of add outcome: " + err.Error())
		return err
	}

	response, err := service.ChangeOutcome(&request)
	if err != nil {
		fmt.Println("An error occurred during getting break " + string(body) + ": " + err.Error())
		return err
	}

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("An error occurred during getting marshalling break response: " + err.Error())
		return err
	}

	w.Write(data)

	return nil
}
