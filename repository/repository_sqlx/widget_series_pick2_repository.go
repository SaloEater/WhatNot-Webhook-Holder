package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type WidgetSeriesPick2Repository struct {
	DB *sqlx.DB
}

func (r *WidgetSeriesPick2Repository) GetByChannel(channelId int64) (*entity.WidgetSeriesPick2, error) {
	var w entity.WidgetSeriesPick2
	err := r.DB.Unsafe().Get(&w, `SELECT * FROM widget_series_pick2 WHERE channel_id = $1`, channelId)
	if err != nil {
		return &entity.WidgetSeriesPick2{ChannelId: channelId, Price: 0}, nil
	}
	return &w, nil
}

func (r *WidgetSeriesPick2Repository) Upsert(w *entity.WidgetSeriesPick2) error {
	_, err := r.DB.NamedExec(`
		INSERT INTO widget_series_pick2 (channel_id, price)
		VALUES (:channel_id, :price)
		ON CONFLICT (channel_id) DO UPDATE SET price = :price`, w)
	return err
}
