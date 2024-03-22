package service

import (
	"encoding/json"
	"errors"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"os"
)

type AddBreakRequest struct {
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	Day       int    `json:"day"`
	Name      string `json:"name"`
	StartDate int64  `json:"start_date"`
	EndDate   int64  `json:"end_date"`
}

func AddBreak(r *AddBreakRequest) error {
	daysData, err := GetDays()
	if err != nil {
		return err
	}

	var days entity.Days
	err = json.Unmarshal(daysData, &days)
	if err != nil {
		return err
	}

	dayIsMissing := true
	breakWithSameNameExists := false
	for i, day := range days.Days {
		if day.Date.Day == r.Day && day.Date.Month == r.Month && day.Date.Year == r.Year {
			dayIsMissing = false
			for _, breakName := range day.Breaks {
				if breakName == r.Name {
					breakWithSameNameExists = true
					break
				}
			}
			dayP := &days.Days[i]
			dayP.Breaks = append(day.Breaks, r.Name)
			break
		}
	}

	if dayIsMissing {
		return errors.New("day is not found")
	}

	if breakWithSameNameExists {
		return errors.New("break with the same name already exists")
	}

	breakFilepath := getFilepath(dataDir, createBreakFilename(r.Year, r.Month, r.Day, r.Name))
	emptyBreak := entity.Break{
		Name:      r.Name,
		Events:    []entity.Event{},
		StartDate: r.StartDate,
		EndDate:   r.EndDate,
	}
	emptyBreakData, err := json.Marshal(emptyBreak)
	if err != nil {
		return err
	}

	err = os.WriteFile(breakFilepath, emptyBreakData, 0644)
	if err != nil {
		return err
	}

	return updateDaysFile(days)
}
