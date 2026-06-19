package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/cache"

type SeriesUpdateRequest struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	UsedCards    int64  `json:"used_cards"`
	DefaultPrice string `json:"default_price"`
}

type SeriesUpdateResponse struct {
	Success bool `json:"success"`
}

func (s *Service) SeriesUpdate(r *SeriesUpdateRequest) (*SeriesUpdateResponse, error) {
	response := &SeriesUpdateResponse{Success: false}
	err := s.SeriesRepositorier.Update(r.Id, r.Name, r.UsedCards, r.DefaultPrice)
	if err == nil {
		response.Success = true
		s.SeriesWithCountCache.Delete(cache.IdToKey(r.Id))
	}
	return response, err
}
