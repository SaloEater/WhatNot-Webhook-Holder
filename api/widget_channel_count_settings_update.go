package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
)

func (a *API) UpdateWidgetChannelCountSettings(w http.ResponseWriter, r *http.Request) (any, error) {
	request := service.UpdateWidgetChannelCountSettingsRequest{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of update widget channel count settings: " + err.Error())
		return nil, err
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of update widget channel count settings: " + err.Error())
		return nil, err
	}
	return a.Service.UpdateWidgetChannelCountSettings(&request)
}
