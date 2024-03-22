package service

import (
	"encoding/json"
	"errors"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"os"
)

type AddDayRequest struct {
	Year  int
	Month int
	Day   int
}

func AddDay(r *AddDayRequest) error {
	daysData, err := GetDays()
	if err != nil {
		return err
	}

	var days entity.Days
	err = json.Unmarshal(daysData, &days)
	if err != nil {
		return err
	}

	found := false
	for _, day := range days.Days {
		if day.Date.Day == r.Day && day.Date.Month == r.Month && day.Date.Year == r.Year {
			found = true
			break
		}
	}

	if found {
		return errors.New("day already exists")
	}

	newDay := entity.Day{
		Date: entity.Date{
			Day:   r.Day,
			Month: r.Month,
			Year:  r.Year,
		},
		Breaks: []string{},
	}
	days.Days = append(days.Days, newDay)

	return updateDaysFile(days)
}

func updateDaysFile(days entity.Days) error {
	daysPath := getFilepath(dataDir, daysFile)
	daysDataUpdate, err := json.Marshal(days)
	if err != nil {
		return err
	}

	err = os.WriteFile(daysPath, daysDataUpdate, 0644)
	if err != nil {
		return err
	}

	return nil
}
