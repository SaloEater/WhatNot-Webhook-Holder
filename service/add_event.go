package service

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"os"
)

type AddEventRequest struct {
	Year       int
	Month      int
	Day        int
	Name       string `json:"name"`
	Customer   string
	Price      float32
	Team       string
	IsGiveaway bool `json:"is_giveaway"`
	Note       string
	Quantity   int
}

func AddEvent(r *AddEventRequest) error {
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

	dayBreak.Events = append(dayBreak.Events, entity.Event{
		Id: RandStringBytes(5),

		Customer:   r.Customer,
		Price:      r.Price,
		Team:       r.Customer,
		IsGiveaway: r.IsGiveaway,
		Note:       r.Customer,
		Quantity:   r.Quantity,
	})
	data, err := json.Marshal(dayBreak)
	err = os.WriteFile(breakFilepath, data, 0644)
	if err != nil {
		fmt.Println("An error occurred during writing writing")
	}

	return nil
}
