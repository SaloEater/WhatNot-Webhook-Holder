package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type WidgetCardsBoardSettingsRepositorier interface {
	GetByChannel(channelId int64) (*entity.WidgetCardsBoardSettings, error)
	Upsert(w *entity.WidgetCardsBoardSettings) error
}
