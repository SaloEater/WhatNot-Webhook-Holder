package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetEventsByBreakRequest struct {
	BreakId int64 `json:"break_id"`
}

type GetEventsByBreakResponse struct {
	Events []*entity.Event `json:"events"`
}

func (s *Service) GetEventsByBreak(r *GetEventsByBreakRequest) (*GetEventsByBreakResponse, error) {
	events, err := s.EventRepository.GetAllByBreak(r.BreakId)
	if err != nil {
		return nil, err
	}

	return &GetEventsByBreakResponse{
		Events: events,
	}, nil
}
