package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type WidgetSeriesBoxesPerBreakRepository struct {
	DB *sqlx.DB
}

func (r *WidgetSeriesBoxesPerBreakRepository) GetBySeries(seriesId int64) (*entity.WidgetSeriesBoxesPerBreak, error) {
	var w entity.WidgetSeriesBoxesPerBreak
	err := r.DB.Get(&w, `SELECT * FROM widget_series_boxes_per_break WHERE series_id = $1`, seriesId)
	if err != nil {
		return &entity.WidgetSeriesBoxesPerBreak{SeriesId: seriesId, Amount: 0}, nil
	}
	return &w, nil
}

func (r *WidgetSeriesBoxesPerBreakRepository) Upsert(w *entity.WidgetSeriesBoxesPerBreak) error {
	_, err := r.DB.NamedExec(`
		INSERT INTO widget_series_boxes_per_break (series_id, amount)
		VALUES (:series_id, :amount)
		ON CONFLICT (series_id) DO UPDATE SET amount = :amount`, w)
	return err
}
