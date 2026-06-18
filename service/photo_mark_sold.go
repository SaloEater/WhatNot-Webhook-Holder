package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/cache"

type PhotoMarkSoldRequest struct {
	Id       int64 `json:"id"`
	SeriesId int64 `json:"series_id"`
	Sold     bool  `json:"sold"`
}

type PhotoMarkSoldResponse struct {
	Success bool `json:"success"`
}

func (s *Service) PhotoMarkSold(r *PhotoMarkSoldRequest) (*PhotoMarkSoldResponse, error) {
	response := &PhotoMarkSoldResponse{Success: false}
	err := s.PhotoRepositorier.MarkSold(r.Id, r.Sold)
	if err == nil {
		response.Success = true
		s.SeriesPricesCache.Delete(cache.IdToKey(r.SeriesId))
	}
	return response, err
}
