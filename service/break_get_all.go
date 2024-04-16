package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetBreaksByStreamRequest struct {
	DayId int64 `json:"day_id"`
}

func (s *Service) GetBreaksByDay(r *GetBreaksByStreamRequest) ([]*entity.Break, error) {
	breaks, err := s.BreakRepository.GetBreaksByDay(r.DayId)
	if err != nil {
		return nil, err
	}

	return breaks, nil
}
