package repository_sqlx

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type DayRepository struct {
	DB *sqlx.DB
}

func (r *DayRepository) GetAll() ([]*entity.Stream, error) {
	days := []*entity.Stream{}
	err := r.DB.Select(&days, `SELECT * FROM stream WHERE is_deleted = false`)
	return days, err
}

func (r *DayRepository) Get(id int64) (*entity.Stream, error) {
	var day entity.Stream
	err := r.DB.Get(&day, `SELECT * FROM stream where id = $1 AND is_deleted = false`, id)
	return &day, err
}

func (r *DayRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`
WITH updatedStream AS (UPDATE stream SET is_deleted = TRUE WHERE id = $1),
breakIDs AS (SELECT id FROM break WHERE day_id=$1),
updatedBreaks AS (UPDATE break SET is_deleted = true WHERE id IN (SELECT * FROM breakIDs)),
updatedEvents AS (UPDATE event SET is_deleted = true WHERE break_id IN (SELECT * FROM breakIDs))
SELECT TRUE

`, id)
	return err
}

func (r *DayRepository) Update(day *entity.Stream) error {
	_, err := r.DB.NamedExec(`UPDATE stream SET
		created_at = :created_at
	WHERE id = :id`, day)

	return err
}

func (r *DayRepository) Create(day *entity.Stream) (int64, error) {
	var id int64
	rows, err := r.DB.NamedQuery(`INSERT INTO stream (
		name, created_at, is_deleted
	) VALUES (
	  	:name,
		:created_at,
		:is_deleted
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

func (r *DayRepository) GetUsernames(id int64) ([]string, error) {
	usernames := []string{}
	err := r.DB.Select(&usernames, `
SELECT DISTINCT customer from event
WHERE is_deleted IS FALSE
GROUP BY customer
HAVING COUNT(*) > 10
`)

	return usernames, err
}
