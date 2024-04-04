package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetDaysResponse struct {
	Days []*entity.Day `json:"days"`
}

func (s *Service) GetDays() (*GetDaysResponse, error) {
	days, err := s.DayRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return &GetDaysResponse{Days: days}, nil
}
