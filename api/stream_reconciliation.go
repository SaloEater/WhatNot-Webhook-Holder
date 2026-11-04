package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
)

func (a *API) StreamReconciliation(w http.ResponseWriter, r *http.Request) (any, error) {
	request := service.StreamReconciliationRequest{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of stream reconciliation: " + err.Error())
		return nil, err
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of stream reconciliation: " + err.Error())
		return nil, err
	}
	return a.Service.StreamReconciliation(&request)
}

func (a *API) ReconciliationSaveManualMap(w http.ResponseWriter, r *http.Request) (any, error) {
	request := service.ReconciliationSaveManualMapRequest{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading body of reconciliation manual map: " + err.Error())
		return nil, err
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshalling body of reconciliation manual map: " + err.Error())
		return nil, err
	}
	return a.Service.ReconciliationSaveManualMap(&request)
}
