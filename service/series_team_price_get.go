package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type SeriesTeamPriceGetRequest struct {
	SeriesId int64 `json:"series_id"`
}

func (s *Service) SeriesTeamPriceGet(r *SeriesTeamPriceGetRequest) ([]*entity.SeriesTeamPrice, error) {
	return s.SeriesTeamPriceRepositorier.GetBySeriesId(r.SeriesId)
}
