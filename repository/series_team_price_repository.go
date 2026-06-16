package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type SeriesTeamPriceRepositorier interface {
	Set(seriesId int64, team string, price float64) error
	GetBySeriesId(seriesId int64) ([]*entity.SeriesTeamPrice, error)
	GetLastPrices() ([]float64, error)
}
