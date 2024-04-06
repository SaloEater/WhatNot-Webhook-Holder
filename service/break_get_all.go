package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetBreaksByDayRequest struct {
	DayId int64 `json:"day_id"`
}

type GetBreaksByDayResponse struct {
	Breaks []*entity.Break `json:"breaks"`
}

func (s *Service) GetBreaksByDay(r *GetBreaksByDayRequest) (*GetBreaksByDayResponse, error) {
	breaks, err := s.BreakRepository.GetBreaksByDay(r.DayId)
	if err != nil {
		return nil, err
	}

	return &GetBreaksByDayResponse{
		Breaks: breaks,
	}, nil
}
