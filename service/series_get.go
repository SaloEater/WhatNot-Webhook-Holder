package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type SeriesGetRequest struct {
	Id int64 `json:"id"`
}

func (s *Service) SeriesGet(r *SeriesGetRequest) (*entity.Series, error) {
	return s.SeriesRepositorier.Get(r.Id)
}
