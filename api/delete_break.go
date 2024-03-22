package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
)

func DeleteBreak(w http.ResponseWriter, r *http.Request) error {
	request := service.DeleteBreakRequest{}
	var err error
	var body []byte

	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of delete break: " + err.Error())
		return err
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of delete break: " + err.Error())
		return err
	}

	err = service.DeleteBreak(&request)
	if err != nil {
		fmt.Println("An error occurred during deleting break " + string(body) + ": " + err.Error())
		return err
	}

	return nil
}
