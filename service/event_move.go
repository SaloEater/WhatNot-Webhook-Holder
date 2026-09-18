package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/cache"

type MoveEventRequest struct {
	Id       int64
	NewIndex int `json:"new_index"`
}

type MoveEventResponse struct {
	Success bool `json:"success"`
}

func (s *Service) MoveEvent(r *MoveEventRequest) (*MoveEventResponse, error) {
	response := &MoveEventResponse{}
	// Request only carries the event id, not the break id, so fetch the event first to learn it.
	event, err := s.EventRepositorier.Get(r.Id)
	if err != nil {
		return response, err
	}

	err = s.EventRepositorier.Move(r.Id, r.NewIndex)
	if err == nil {
		response.Success = true
		s.EventsCache.Delete(cache.IdToKey(event.BreakId))
	}

	return response, err
}
