package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetBreakEventsRequest struct {
	BreakId int64 `json:"break_id"`
}

type GetBreakEventsResponse struct {
	Events []*entity.Event `json:"events"`
}

func (s *Service) GetBreakEvents(r *GetBreakEventsRequest) (*GetBreakEventsResponse, error) {
	key := cache.IdToKey(r.BreakId)

	if cached, found := s.EventsCache.Get(key); found {
		return &GetBreakEventsResponse{
			Events: cached,
		}, nil
	}

	events, err := s.EventRepositorier.GetAllByBreak(r.BreakId)
	if err != nil {
		return nil, err
	}
	s.EventsCache.Set(key, events)

	return &GetBreakEventsResponse{
		Events: events,
	}, nil
}
