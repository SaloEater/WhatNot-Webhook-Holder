package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type SeriesGetPricesRequest struct {
	SeriesId int64 `json:"series_id"`
}

func (s *Service) SeriesGetPrices(r *SeriesGetPricesRequest) ([]*entity.SeriesTeamTotal, error) {
	key := cache.IdToKey(r.SeriesId)
	if s.SeriesPricesCache.Has(key) {
		cached, _ := s.SeriesPricesCache.Get(key)
		return cached, nil
	}
	totals, err := s.PhotoRepositorier.GetPricesBySeriesId(r.SeriesId)
	if err == nil {
		s.SeriesPricesCache.Set(key, totals)
	}
	return totals, err
}
