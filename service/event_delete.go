package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/cache"

const notFound = -1

type DeleteEventRequest struct {
	Id int64 `json:"id"`
}

type DeleteEventResponse struct {
	Success bool `json:"success"`
}

func (s *Service) DeleteEvent(r *DeleteEventRequest) (*DeleteEventResponse, error) {
	response := &DeleteEventResponse{Success: false}
	// Request only carries the event id, not the break id, so fetch the event first to learn it.
	event, err := s.EventRepositorier.Get(r.Id)
	if err != nil {
		return response, err
	}

	err = s.EventRepositorier.Delete(r.Id)
	if err == nil {
		response.Success = true
		s.EventsCache.Delete(cache.IdToKey(event.BreakId))
	}

	return response, err
}
