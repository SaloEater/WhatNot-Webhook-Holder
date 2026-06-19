package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type WidgetSeriesPick2Repositorier interface {
	GetByChannel(channelId int64) (*entity.WidgetSeriesPick2, error)
	Upsert(w *entity.WidgetSeriesPick2) error
}
