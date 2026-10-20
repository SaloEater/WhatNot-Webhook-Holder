package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type SeriesPresetRepositorier interface {
	ListByChannel(channelId int64) ([]*entity.SeriesPreset, error)
	Create(channelId int64, name string, seriesIds []int64) (int64, error)
	Delete(id int64) error
}
