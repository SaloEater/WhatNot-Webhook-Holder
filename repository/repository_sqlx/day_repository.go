package repository_sqlx

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type DayRepository struct {
	DB *sqlx.DB
}

func (r *DayRepository) GetAll() ([]*entity.Day, error) {
	days := []*entity.Day{}
	err := r.DB.Select(&days, `SELECT * FROM day`)
	return days, err
}

func (r *DayRepository) Get(id int64) (*entity.Day, error) {
	var day entity.Day
	err := r.DB.Get(&day, `SELECT * FROM day where id = $1`, id)
	return &day, err
}

func (r *DayRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM day WHERE id = $1`, id)
	return err
}

func (r *DayRepository) Update(day *entity.Day) error {
	_, err := r.DB.NamedExec(`UPDATE day SET
		date = :date
	WHERE id = :id`, day)

	return err
}

func (r *DayRepository) Create(day *entity.Day) (int64, error) {
	var id int64
	rows, err := r.DB.NamedQuery(`INSERT INTO day (
		date
	) VALUES (
		:date
	) RETURNING (id)`, day)
	if err != nil {
		return id, err
	}

	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("no rows returned after INSERT")
	}
	return id, err
}
