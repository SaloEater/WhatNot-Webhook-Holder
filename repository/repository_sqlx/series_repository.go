package repository_sqlx

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type SeriesRepository struct {
	DB *sqlx.DB
}

func (r *SeriesRepository) Create(name string) (int64, error) {
	var id int64
	err := r.DB.QueryRow(`INSERT INTO series (name) VALUES ($1) RETURNING id`, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SeriesRepository) Get(id int64) (*entity.Series, error) {
	var s entity.Series
	err := r.DB.Get(&s, `SELECT * FROM series WHERE id = $1 AND is_deleted = false`, id)
	return &s, err
}

func (r *SeriesRepository) GetList() ([]*entity.Series, error) {
	series := []*entity.Series{}
	err := r.DB.Select(&series, `SELECT * FROM series WHERE is_deleted = false ORDER BY created_at DESC`)
	return series, err
}

func (r *SeriesRepository) Update(id int64, name string) error {
	_, err := r.DB.Exec(`UPDATE series SET name = $1 WHERE id = $2`, name, id)
	return err
}

func (r *SeriesRepository) Close(id int64) error {
	_, err := r.DB.Exec(`UPDATE series SET status = 'closed' WHERE id = $1`, id)
	return err
}

func (r *SeriesRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`UPDATE series SET is_deleted = true WHERE id = $1`, id)
	return err
}

func (r *SeriesRepository) CountOpen() (int, error) {
	var count int
	err := r.DB.QueryRow(`SELECT COUNT(*) FROM series WHERE status = 'open' AND is_deleted = false`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("CountOpen: %w", err)
	}
	return count, nil
}
