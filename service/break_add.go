package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"strconv"
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
	startDate, err := strconv.ParseInt(r.StartDate, 10, 0)
	if err != nil {
		return nil, err
	}
	endDate, err := strconv.ParseInt(r.EndDate, 10, 0)
	if err != nil {
		return nil, err
	}

	var id int64
	id, err = s.BreakRepository.Create(&entity.Break{
		DayId:     r.DayId,
		Name:      r.Name,
		StartDate: time.UnixMilli(startDate),
		EndDate:   time.UnixMilli(endDate),
	})
	if err != nil {
		return nil, err
	}

	return &AddBreakResponse{Id: id}, nil

}
