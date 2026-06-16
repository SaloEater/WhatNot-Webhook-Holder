package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type SeriesTeamPriceRepository struct {
	DB *sqlx.DB
}

func (r *SeriesTeamPriceRepository) Set(seriesId int64, team string, price float64) error {
	_, err := r.DB.Exec(`
INSERT INTO series_team_price (series_id, team, price)
VALUES ($1, $2, $3)
ON CONFLICT (series_id, team) DO UPDATE SET price = EXCLUDED.price
`, seriesId, team, price)
	return err
}

func (r *SeriesTeamPriceRepository) GetBySeriesId(seriesId int64) ([]*entity.SeriesTeamPrice, error) {
	prices := []*entity.SeriesTeamPrice{}
	err := r.DB.Select(&prices, `SELECT * FROM series_team_price WHERE series_id = $1 ORDER BY team`, seriesId)
	return prices, err
}

func (r *SeriesTeamPriceRepository) GetLastPrices() ([]float64, error) {
	prices := []float64{}
	err := r.DB.Select(&prices, `
		SELECT DISTINCT stp.price
		FROM series_team_price stp
		JOIN series s ON s.id = stp.series_id
		WHERE stp.price > 0
		  AND s.is_deleted = false
		ORDER BY stp.price
	`)
	return prices, err
}
