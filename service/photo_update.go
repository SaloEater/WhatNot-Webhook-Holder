package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/cache"

type PhotoUpdateRequest struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Team  string `json:"team"`
	Price int64  `json:"price"`
}

type PhotoUpdateResponse struct {
	Success bool `json:"success"`
}

func (s *Service) PhotoUpdate(r *PhotoUpdateRequest) (*PhotoUpdateResponse, error) {
	response := &PhotoUpdateResponse{Success: false}

	photo, err := s.PhotoRepositorier.GetById(r.Id)
	if err != nil {
		return response, err
	}
	priceChanged := photo.Price != r.Price

	err = s.PhotoRepositorier.Update(r.Id, r.Name, r.Team, r.Price)
	if err == nil {
		response.Success = true
		if priceChanged {
			s.SeriesPricesCache.Delete(cache.IdToKey(photo.SeriesId))
		}
	}
	return response, err
}
