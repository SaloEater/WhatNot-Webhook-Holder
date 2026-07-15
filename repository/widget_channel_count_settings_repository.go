package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type WidgetChannelCountSettingsRepositorier interface {
	GetByChannel(channelId int64) (*entity.WidgetChannelCountSettings, error)
	Upsert(w *entity.WidgetChannelCountSettings) error
}
