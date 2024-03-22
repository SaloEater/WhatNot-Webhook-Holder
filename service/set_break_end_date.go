package service

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"os"
	"strconv"
)

type SetBreakEndDateRequest struct {
	Year    int
	Month   int
	Day     int
	Name    string
	EndDate string `json:"end_date"`
}

func SetBreakEndDate(r *SetBreakEndDateRequest) error {
	breakFilepath := getFilepath(dataDir, createBreakFilename(r.Year, r.Month, r.Day, r.Name))
	breakData, err := os.ReadFile(breakFilepath)
	if err != nil {
		return err
	}

	var dayBreak entity.Break
	err = json.Unmarshal(breakData, &dayBreak)
	if err != nil {
		return err
	}

	dayBreak.EndDate, err = strconv.ParseInt(r.EndDate, 10, 0)
	if err != nil {
		return err
	}

	data, err := json.Marshal(dayBreak)
	if err != nil {
		return err
	}

	err = os.WriteFile(breakFilepath, data, 0644)
	if err != nil {
		fmt.Println("An error occurred during writing writing")
	}

	return err
}
