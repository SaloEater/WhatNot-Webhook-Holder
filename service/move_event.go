package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"os"
)

type MoveEventRequest struct {
	Year      int
	Month     int
	Day       int
	BreakName string `json:"break_name"`
	Id        string
	NewIndex  int `json:"new_index"`
}

func MoveEvent(r *MoveEventRequest) error {

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

	oldEventIndex := notFound
	for i, event := range dayBreak.Events {
		if event.Id == r.Id {
			oldEventIndex = i
		}
	}

	if oldEventIndex == notFound {
		return errors.New("break is not found")
	}

	if oldEventIndex == r.NewIndex {
		return nil
	}

	if oldEventIndex < 0 || r.NewIndex >= len(dayBreak.Events) {
		return errors.New("index is invalid")
	}

	var newEvents []entity.Event
	if oldEventIndex < r.NewIndex {
		newEvents = append(dayBreak.Events[:oldEventIndex], dayBreak.Events[oldEventIndex+1:r.NewIndex+1]...) // elements before new position without element
		newEvents = append(newEvents, dayBreak.Events[oldEventIndex])                                         // element
		newEvents = append(newEvents, dayBreak.Events[r.NewIndex+2:]...)                                      // elements after new position
	} else {
		newEvents = append(dayBreak.Events[:r.NewIndex], dayBreak.Events[oldEventIndex]) // elements before new position without element
		newEvents = append(newEvents, dayBreak.Events[oldEventIndex:r.NewIndex]...)      // element
		newEvents = append(newEvents, dayBreak.Events[oldEventIndex+1:]...)              // elements after new position
	}

	dayBreak.Events = newEvents

	data, err := json.Marshal(dayBreak)
	err = os.WriteFile(breakFilepath, data, 0644)
	if err != nil {
		fmt.Println("An error occurred during writing writing")
	}

	return nil
}
