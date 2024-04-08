package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetBreakRequest struct {
	Id int64 `json:"id"`
}

func (s *Service) GetBreak(r *GetBreakRequest) (*entity.Break, error) {
	dayBreak, err := s.BreakRepository.Get(r.Id)
	if err != nil {
		return nil, err
	}

	return dayBreak, nil
}
