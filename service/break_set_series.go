package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/cache"

type BreakSetSeriesRequest struct {
	BreakId  int64 `json:"break_id"`
	SeriesId int64 `json:"series_id"`
}

type BreakSetSeriesResponse struct {
	Success bool `json:"success"`
}

func (s *Service) BreakSetSeries(r *BreakSetSeriesRequest) (*BreakSetSeriesResponse, error) {
	response := &BreakSetSeriesResponse{Success: false}
	err := s.BreakRepositorier.SetSeries(r.BreakId, &r.SeriesId)
	if err == nil {
		response.Success = true
		s.BreakCache.Delete(cache.IdToKey(r.BreakId))
	}
	return response, err
}
