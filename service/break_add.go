package service

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"time"
)

type AddBreakRequest struct {
	DayId     int64  `json:"day_id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type AddBreakResponse struct {
	Id int64 `json:"id"`
}

func (s *Service) AddBreak(r *AddBreakRequest) (*AddBreakResponse, error) {
	startDate, err := time.Parse(DateLayout, r.StartDate)
	if err != nil {
		fmt.Println("Error parsing start date string:", err)
		return nil, err
	}

	endDate, err := time.Parse(DateLayout, r.EndDate)
	if err != nil {
		fmt.Println("Error parsing end date string:", err)
		return nil, err
	}

	channel, err := s.ChannelRepository.GetByStream(r.DayId)
	if err != nil {
		return nil, err
	}

	var id int64
	id, err = s.BreakRepository.Create(&entity.Break{
		DayId:        r.DayId,
		Name:         r.Name,
		StartDate:    startDate,
		EndDate:      endDate,
		HighBidTeam:  channel.DefaultHighBidTeam,
		GiveawayTeam: "",
		IsDeleted:    false,
		HighBidFloor: channel.DefaultHighBidFloor,
	})
	if err != nil {
		return nil, err
	}

	return &AddBreakResponse{Id: id}, nil

}
