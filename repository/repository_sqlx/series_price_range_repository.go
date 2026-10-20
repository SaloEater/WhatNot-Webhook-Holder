package repository_sqlx

import (
	"errors"
	"fmt"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type SeriesPriceRangeRepository struct {
	DB *sqlx.DB
}

func (r *SeriesPriceRangeRepository) ListBySeries(seriesId int64) ([]*entity.SeriesPriceRange, error) {
	rows := []*entity.SeriesPriceRange{}
	err := r.DB.Unsafe().Select(&rows, `SELECT * FROM series_price_ranges WHERE series_id = $1 ORDER BY sort_order ASC, id ASC`, seriesId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *SeriesPriceRangeRepository) Create(w *entity.SeriesPriceRange) (int64, error) {
	var id int64
	err := r.DB.QueryRow(
		`INSERT INTO series_price_ranges (series_id, price_from, price_to, count, label, sort_order)
		VALUES ($1, $2, $3, $4, $5, (SELECT COALESCE(MAX(sort_order), -1) + 1 FROM series_price_ranges WHERE series_id = $1)) RETURNING id`,
		w.SeriesId, w.PriceFrom, w.PriceTo, w.Count, w.Label,
	).Scan(&id)
	return id, err
}

func (r *SeriesPriceRangeRepository) Update(w *entity.SeriesPriceRange) error {
	_, err := r.DB.Exec(
		`UPDATE series_price_ranges SET price_from = $1, price_to = $2, count = $3, label = $4 WHERE id = $5`,
		w.PriceFrom, w.PriceTo, w.Count, w.Label, w.Id,
	)
	return err
}

func (r *SeriesPriceRangeRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM series_price_ranges WHERE id=$1`, id)
	return err
}

func (r *SeriesPriceRangeRepository) Reorder(seriesId int64, ids []int64) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	existing := []int64{}
	if err := tx.Select(&existing, `SELECT id FROM series_price_ranges WHERE series_id = $1 FOR UPDATE`, seriesId); err != nil {
		return err
	}
	if len(ids) != len(existing) {
		return errors.New("ids must contain every price range of the series exactly once")
	}
	belongs := make(map[int64]bool, len(existing))
	for _, id := range existing {
		belongs[id] = true
	}
	seen := make(map[int64]bool, len(ids))
	for _, id := range ids {
		if !belongs[id] {
			return fmt.Errorf("price range %d does not belong to series %d", id, seriesId)
		}
		if seen[id] {
			return fmt.Errorf("price range %d is listed more than once", id)
		}
		seen[id] = true
	}
	for i, id := range ids {
		if _, err := tx.Exec(`UPDATE series_price_ranges SET sort_order = $1 WHERE id = $2`, i, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}
