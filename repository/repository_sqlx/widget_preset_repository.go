package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type WidgetPresetRepository struct {
	DB *sqlx.DB
}

func (r *WidgetPresetRepository) ListByChannel(channelId int64) ([]*entity.WidgetPreset, error) {
	var rows []*entity.WidgetPreset
	err := r.DB.Unsafe().Select(&rows, `SELECT * FROM widget_presets WHERE channel_id = $1 ORDER BY name`, channelId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *WidgetPresetRepository) Upsert(channelId int64, name, value string) error {
	_, err := r.DB.Exec(`
		INSERT INTO widget_presets (channel_id, name, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (channel_id, name) DO UPDATE SET value = $3`,
		channelId, name, value)
	return err
}

func (r *WidgetPresetRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM widget_presets WHERE id=$1`, id)
	return err
}
