package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type UpdateEventRequest struct {
	Id           int64   `json:"id"`
	BreakId      int64   `json:"break_id"`
	Customer     string  `json:"customer"`
	Price        float32 `json:"price"`
	Team         string  `json:"team"`
	IsGiveaway   bool    `json:"is_giveaway"`
	Note         string  `json:"note"`
	Quantity     int     `json:"quantity"`
	GiveawayType int16   `json:"giveaway_type"`
}

type UpdateEventResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateEvent(r *UpdateEventRequest) (*UpdateEventResponse, error) {
	response := &UpdateEventResponse{}
	event, err := s.EventRepositorier.Get(r.Id)
	if err != nil {
		return response, err
	}
	// The request may carry a different break_id (the repository's UPDATE writes it), which moves
	// the event between breaks — so BOTH breaks' cached lists go stale, not just the target's.
	prevBreakId := event.BreakId
	event.BreakId = r.BreakId
	event.Customer = r.Customer
	event.Price = r.Price
	event.Team = r.Team
	event.IsGiveaway = r.IsGiveaway
	event.Note = r.Note
	event.Quantity = r.Quantity
	event.GiveawayType = entity.GiveawayType(r.GiveawayType)
	err = s.EventRepositorier.Update(event)
	if err == nil {
		response.Success = true
		s.EventsCache.Delete(cache.IdToKey(r.BreakId))
		if prevBreakId != r.BreakId {
			s.EventsCache.Delete(cache.IdToKey(prevBreakId))
		}
	}

	return response, nil
}
