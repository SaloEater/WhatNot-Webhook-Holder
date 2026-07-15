package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type WidgetChannelCountSettingsRepository struct {
	DB *sqlx.DB
}

func (r *WidgetChannelCountSettingsRepository) GetByChannel(channelId int64) (*entity.WidgetChannelCountSettings, error) {
	var w entity.WidgetChannelCountSettings
	err := r.DB.Get(&w, `SELECT * FROM widget_channel_count_settings WHERE channel_id = $1`, channelId)
	if err != nil {
		def := &entity.WidgetChannelCountSettings{ChannelId: channelId, ShowPercentage: false}
		_ = r.Upsert(def)
		return def, nil
	}
	return &w, nil
}

func (r *WidgetChannelCountSettingsRepository) Upsert(w *entity.WidgetChannelCountSettings) error {
	_, err := r.DB.NamedExec(`
		INSERT INTO widget_channel_count_settings (channel_id, show_percentage)
		VALUES (:channel_id, :show_percentage)
		ON CONFLICT (channel_id) DO UPDATE SET show_percentage = :show_percentage`, w)
	return err
}
