package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type UpdateAllEventsRequest struct {
	Events []*UpdateEventRequest `json:"events"`
}

type UpdateAllEventsResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateAllEvents(r *UpdateAllEventsRequest) (*UpdateAllEventsResponse, error) {
	response := &UpdateAllEventsResponse{}
	ids := make([]int64, len(r.Events))

	eventsMap := make(map[int64]*UpdateEventRequest, len(ids))
	for i, event := range r.Events {
		ids[i] = event.Id
		eventsMap[event.Id] = event
	}

	events, err := s.EventRepositorier.GetAll(ids)
	if err != nil {
		return response, err
	}

	for i := range events {
		event := events[i]
		requestEvent := eventsMap[event.Id]
		event.BreakId = requestEvent.BreakId
		event.Customer = requestEvent.Customer
		event.Price = requestEvent.Price
		event.Team = requestEvent.Team
		event.IsGiveaway = requestEvent.IsGiveaway
		event.Note = requestEvent.Note
		event.Quantity = requestEvent.Quantity
		event.GiveawayType = entity.GiveawayType(requestEvent.GiveawayType)
		events[i] = event
	}

	err = s.EventRepositorier.UpdateAll(events)
	if err == nil {
		response.Success = true
		// The batch may span multiple breaks, so invalidate every distinct break id in it.
		seenBreakIds := make(map[int64]bool, len(events))
		for _, event := range events {
			if !seenBreakIds[event.BreakId] {
				seenBreakIds[event.BreakId] = true
				s.EventsCache.Delete(cache.IdToKey(event.BreakId))
			}
		}
	}

	return response, nil
}
