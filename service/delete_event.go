package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"os"
)

const notFound = -1

type DeleteEventRequest struct {
	Year      int
	Month     int
	Day       int
	BreakName string
	Id        string
}

func DeleteEvent(r *DeleteEventRequest) error {
	breakFilepath := getFilepath(dataDir, createBreakFilename(r.Year, r.Month, r.Day, r.BreakName))
	breakData, err := os.ReadFile(breakFilepath)
	if err != nil {
		return err
	}

	var dayBreak entity.Break
	err = json.Unmarshal(breakData, &dayBreak)
	if err != nil {
		return err
	}

	eventIndex := notFound
	for i, event := range dayBreak.Events {
		if event.Id == r.Id {
			eventIndex = i
		}
	}

	if eventIndex == notFound {
		return errors.New("break is not found")
	}

	dayBreak.Events = append(dayBreak.Events[:eventIndex], dayBreak.Events[eventIndex+1:]...)

	data, err := json.Marshal(dayBreak)
	err = os.WriteFile(breakFilepath, data, 0644)
	if err != nil {
		fmt.Println("An error occurred during writing writing")
	}

	return nil
}
