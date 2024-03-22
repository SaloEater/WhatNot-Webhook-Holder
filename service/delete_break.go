package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"os"
)

type DeleteBreakRequest struct {
	Year  int
	Month int
	Day   int
	Name  string
}

func DeleteBreak(r *DeleteBreakRequest) error {
	daysData, err := GetDays()
	if err != nil {
		return err
	}

	var days entity.Days
	err = json.Unmarshal(daysData, &days)
	if err != nil {
		return err
	}

	breakIsMissing := true
	for i, day := range days.Days {
		if day.Date.Day == r.Day && day.Date.Month == r.Month && day.Date.Year == r.Year {
			dayP := &days.Days[i]
			breakIndex := -1
			for j, breakName := range day.Breaks {
				if breakName == r.Name {
					breakIsMissing = false
					breakIndex = j
					break
				}
			}

			if breakIndex != -1 {
				dayP.Breaks = append(day.Breaks[:breakIndex], day.Breaks[breakIndex+1:]...)
				breakFilepath := getFilepath(dataDir, createBreakFilename(r.Year, r.Month, r.Day, r.Name))
				deleteBreakFilepath := getFilepath(dataDir, createDeletedBreakFilename(r.Year, r.Month, r.Day, r.Name))
				err = os.Rename(breakFilepath, deleteBreakFilepath)
			}

			break
		}
	}

	if err != nil {
		fmt.Println("An error occurred during deleting break: " + err.Error())
		return err
	}

	if breakIsMissing {
		return errors.New("break is not found")
	}

	return updateDaysFile(days)
}
