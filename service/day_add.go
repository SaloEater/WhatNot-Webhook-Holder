package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"time"
)

type AddDayRequest struct {
	Year  int
	Month int
	Day   int
}

type AddDayResponse struct {
	Id int64 `json:"id"`
}

func (s *Service) AddDay(r *AddDayRequest) (*AddDayResponse, error) {
	id, err := s.DayRepository.Create(&entity.Day{
		Date: time.Date(r.Year, time.Month(r.Month), r.Day, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, err
	}

	return &AddDayResponse{Id: id}, nil
}
