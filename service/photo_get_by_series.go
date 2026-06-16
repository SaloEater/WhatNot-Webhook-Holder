package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type PhotoGetBySeriesRequest struct {
	SeriesId int64 `json:"series_id"`
}

func (s *Service) PhotoGetBySeries(r *PhotoGetBySeriesRequest) ([]*entity.Photo, error) {
	return s.PhotoRepositorier.GetBySeriesId(r.SeriesId)
}
