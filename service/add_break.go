package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"strconv"
)

type AddBreakRequest struct {
	DayId     int    `json:"day_id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type AddBreakResponse struct {
	Id int `json:"id"`
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

	var id int
	id, err = s.BreakRepository.Create(&entity.Break{
		Name:      r.Name,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return nil, err
	}

	return &AddBreakResponse{Id: id}, nil

}
