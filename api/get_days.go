package api

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"net/http"
)

func GetDays(w http.ResponseWriter, r *http.Request) error {
	daysData, err := service.GetDays()
	if err != nil {
		fmt.Println("An error occurred during getting days" + err.Error())
		return err
	}

	w.Write(daysData)
	return nil
}
