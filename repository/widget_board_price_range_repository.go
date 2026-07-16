package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type WidgetBoardPriceRangeRepositorier interface {
	ListByChannel(channelId int64) ([]*entity.WidgetBoardPriceRange, error)
	Create(w *entity.WidgetBoardPriceRange) (int64, error)
	Upsert(channelId int64, tierId string, priceFrom int) error
	Delete(id int64) error
}
