package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type WidgetSeriesBoxesPerBreakRepositorier interface {
	GetBySeries(seriesId int64) (*entity.WidgetSeriesBoxesPerBreak, error)
	Upsert(w *entity.WidgetSeriesBoxesPerBreak) error
}
