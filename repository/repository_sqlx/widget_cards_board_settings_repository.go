package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type WidgetCardsBoardSettingsRepository struct {
	DB *sqlx.DB
}

func (r *WidgetCardsBoardSettingsRepository) GetByChannel(channelId int64) (*entity.WidgetCardsBoardSettings, error) {
	var w entity.WidgetCardsBoardSettings
	err := r.DB.Get(&w, `SELECT * FROM widget_cards_board_settings WHERE channel_id = $1`, channelId)
	if err != nil {
		def := &entity.WidgetCardsBoardSettings{ChannelId: channelId, Orientation: "list"}
		_ = r.Upsert(def)
		return def, nil
	}
	return &w, nil
}

func (r *WidgetCardsBoardSettingsRepository) Upsert(w *entity.WidgetCardsBoardSettings) error {
	_, err := r.DB.NamedExec(`
		INSERT INTO widget_cards_board_settings (channel_id, orientation)
		VALUES (:channel_id, :orientation)
		ON CONFLICT (channel_id) DO UPDATE SET orientation = :orientation`, w)
	return err
}
