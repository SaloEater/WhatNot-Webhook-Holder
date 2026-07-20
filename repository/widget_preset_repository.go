package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type WidgetPresetRepositorier interface {
	ListByChannel(channelId int64) ([]*entity.WidgetPreset, error)
	Upsert(channelId int64, name, value string) error
	Delete(id int64) error
}
