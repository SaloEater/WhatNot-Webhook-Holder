package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetStreamBreaksRequest struct {
	Id int64 `json:"id"`
}

func (s *Service) GetStreamBreaks(r *GetStreamBreaksRequest) ([]*entity.Break, error) {
	breaks, err := s.BreakRepositorier.GetBreaksByStreamId(r.Id)
	if err != nil {
		return nil, err
	}

	return breaks, nil
}
