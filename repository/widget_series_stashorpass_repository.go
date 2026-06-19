package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type WidgetSeriesStashorpassRepositorier interface {
	GetByChannel(channelId int64) (*entity.WidgetSeriesStashorpass, error)
	Upsert(w *entity.WidgetSeriesStashorpass) error
}
