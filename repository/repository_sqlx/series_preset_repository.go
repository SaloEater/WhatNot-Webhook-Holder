package repository_sqlx

import (
	"encoding/json"
	"time"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type SeriesPresetRepository struct {
	DB *sqlx.DB
}

type seriesPresetRow struct {
	Id        int64     `db:"id"`
	ChannelId int64     `db:"channel_id"`
	Name      string    `db:"name"`
	SeriesIds []byte    `db:"series_ids"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *SeriesPresetRepository) ListByChannel(channelId int64) ([]*entity.SeriesPreset, error) {
	rows := []seriesPresetRow{}
	err := r.DB.Select(&rows, `SELECT id, channel_id, name, series_ids, created_at FROM series_preset WHERE channel_id = $1 ORDER BY name`, channelId)
	if err != nil {
		return nil, err
	}
	presets := make([]*entity.SeriesPreset, 0, len(rows))
	for _, row := range rows {
		ids := []int64{}
		if err := json.Unmarshal(row.SeriesIds, &ids); err != nil {
			return nil, err
		}
		presets = append(presets, &entity.SeriesPreset{
			Id:        row.Id,
			ChannelId: row.ChannelId,
			Name:      row.Name,
			SeriesIds: ids,
			CreatedAt: row.CreatedAt,
		})
	}
	return presets, nil
}

func (r *SeriesPresetRepository) Create(channelId int64, name string, seriesIds []int64) (int64, error) {
	if seriesIds == nil {
		seriesIds = []int64{}
	}
	encoded, err := json.Marshal(seriesIds)
	if err != nil {
		return 0, err
	}
	var id int64
	err = r.DB.QueryRow(
		`INSERT INTO series_preset (channel_id, name, series_ids) VALUES ($1, $2, $3) RETURNING id`,
		channelId, name, string(encoded),
	).Scan(&id)
	return id, err
}

func (r *SeriesPresetRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM series_preset WHERE id = $1`, id)
	return err
}
