package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"os"
)

type UpdateEventRequest struct {
	Year       int
	Month      int
	Day        int
	Name       string `json:"name"`
	Id         string
	Customer   string
	Price      float32
	Team       string
	IsGiveaway bool `json:"is_giveaway"`
	Note       string
	Quantity   int
}

func UpdateEvent(r *UpdateEventRequest) error {
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

	eventIndex := notFound
	for i, event := range dayBreak.Events {
		if event.Id == r.Id {
			eventIndex = i
			pEvent := &dayBreak.Events[i]
			pEvent.Customer = r.Customer
			pEvent.Price = r.Price
			pEvent.Team = r.Team
			pEvent.IsGiveaway = r.IsGiveaway
			pEvent.Note = r.Note
			pEvent.Quantity = r.Quantity
		}
	}

	if eventIndex == notFound {
		return errors.New("break is not found")
	}

	data, err := json.Marshal(dayBreak)
	err = os.WriteFile(breakFilepath, data, 0644)
	if err != nil {
		fmt.Println("An error occurred during writing writing")
	}

	return nil
}
