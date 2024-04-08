package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetDayRequest struct {
	Id int64 `json:"id"`
}

func (s *Service) GetDay(r *GetDayRequest) (*entity.Day, error) {
	day, err := s.DayRepository.Get(r.Id)
	if err != nil {
		return nil, err
	}

	return day, nil
}
