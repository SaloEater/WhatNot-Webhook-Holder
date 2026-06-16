package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
)

func (a *API) PhotoDelete(w http.ResponseWriter, r *http.Request) (any, error) {
	request := service.PhotoDeleteRequest{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of photo delete: " + err.Error())
		return nil, err
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of photo delete: " + err.Error())
		return nil, err
	}
	return a.Service.PhotoDelete(&request)
}
