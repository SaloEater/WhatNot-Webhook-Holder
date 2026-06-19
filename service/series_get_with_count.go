package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type SeriesGetWithCountRequest struct {
	Id int64 `json:"id"`
}

func (s *Service) SeriesGetWithCount(r *SeriesGetWithCountRequest) (*entity.SeriesWithCount, error) {
	key := cache.IdToKey(r.Id)
	if s.SeriesWithCountCache.Has(key) {
		cached, _ := s.SeriesWithCountCache.Get(key)
		return cached, nil
	}
	series, err := s.SeriesRepositorier.Get(r.Id)
	if err != nil {
		return nil, err
	}
	unsoldCount, err := s.PhotoRepositorier.CountUnsoldBySeriesId(r.Id)
	if err != nil {
		return nil, err
	}
	soldCount, err := s.PhotoRepositorier.CountSoldBySeriesId(r.Id)
	if err != nil {
		return nil, err
	}
	result := &entity.SeriesWithCount{Series: series, UnsoldCount: unsoldCount, SoldCount: soldCount}
	s.SeriesWithCountCache.Set(key, result)
	return result, nil
}
