package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"os"
)

type DeleteDayRequest struct {
	Year  int
	Month int
	Day   int
}

func DeleteDay(r *DeleteDayRequest) error {
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
	var dayIndex int
	for i, day := range days.Days {
		if day.Date.Day == r.Day && day.Date.Month == r.Month && day.Date.Year == r.Year {
			found = true
			for _, dayBreakName := range day.Breaks {
				breakFilepath := getFilepath(dataDir, createBreakFilename(r.Year, r.Month, r.Day, dayBreakName))
				deleteBreakFilepath := getFilepath(dataDir, createDeletedBreakFilename(r.Year, r.Month, r.Day, dayBreakName))
				err = os.Rename(breakFilepath, deleteBreakFilepath)
				fmt.Println("An error occurred during deleting break: " + err.Error())
			}
			dayIndex = i
		}
	}

	if !found {
		return errors.New("day is not found")
	}

	days.Days = append(days.Days[:dayIndex], days.Days[dayIndex+1:]...)

	err = updateDaysFile(days)
	if err != nil {
		return err
	}

	return nil
}
