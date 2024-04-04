package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type AddDayRequest struct {
	Year  int
	Month int
	Day   int
}

type AddDayResponse struct {
	Id int
}

func (s *Service) AddDay(r *AddDayRequest) (*AddDayResponse, error) {
	id, err := s.DayRepository.Create(&entity.Day{
		Year:  r.Year,
		Month: r.Month,
		Day:   r.Day,
	})

	if err != nil {
		return nil, err
	}

	return &AddDayResponse{Id: id}, nil
}
