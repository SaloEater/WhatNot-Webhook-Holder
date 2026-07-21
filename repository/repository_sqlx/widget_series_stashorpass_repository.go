package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type WidgetSeriesStashorpassRepository struct {
	DB *sqlx.DB
}

func (r *WidgetSeriesStashorpassRepository) GetByChannel(channelId int64) (*entity.WidgetSeriesStashorpass, error) {
	var w entity.WidgetSeriesStashorpass
	err := r.DB.Unsafe().Get(&w, `SELECT * FROM widget_series_stashorpass WHERE channel_id = $1`, channelId)
	if err != nil {
		return &entity.WidgetSeriesStashorpass{ChannelId: channelId, Price: 0}, nil
	}
	return &w, nil
}

func (r *WidgetSeriesStashorpassRepository) Upsert(w *entity.WidgetSeriesStashorpass) error {
	_, err := r.DB.NamedExec(`
		INSERT INTO widget_series_stashorpass (channel_id, price)
		VALUES (:channel_id, :price)
		ON CONFLICT (channel_id) DO UPDATE SET price = :price`, w)
	return err
}
