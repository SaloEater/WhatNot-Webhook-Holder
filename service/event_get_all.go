package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetBreakEventsRequest struct {
	BreakId int64 `json:"break_id"`
}

type GetBreakEventsResponse struct {
	Events []*entity.Event `json:"events"`
}

func (s *Service) GetBreakEvents(r *GetBreakEventsRequest) (*GetBreakEventsResponse, error) {
	events, err := s.EventRepository.GetAllByBreak(r.BreakId)
	if err != nil {
		return nil, err
	}

	return &GetBreakEventsResponse{
		Events: events,
	}, nil
}
