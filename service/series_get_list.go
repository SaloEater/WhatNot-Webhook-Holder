package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

func (s *Service) SeriesGetList() ([]*entity.Series, error) {
	return s.SeriesRepositorier.GetList()
}
