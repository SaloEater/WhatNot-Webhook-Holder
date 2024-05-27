package repository_sqlx

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type BreakRepository struct {
	DB *sqlx.DB
}

func (r *BreakRepository) Create(dayBreak *entity.Break) (int64, error) {
	var id int64
	rows, err := r.DB.NamedQuery(`INSERT INTO break (
		day_id, name, start_date, end_date, is_deleted, high_bid_team, giveaway_team, high_bid_floor
	) VALUES (
		:day_id,
		:name,
		:start_date,
		:end_date,
		:is_deleted,
		:high_bid_team,
		:giveaway_team,
	  	:high_bid_floor
	) RETURNING (id)`, dayBreak)
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

func (r *BreakRepository) Get(id int64) (*entity.Break, error) {
	var dayBreak entity.Break
	err := r.DB.Get(&dayBreak, `SELECT * FROM break WHERE id = $1 AND is_deleted = false`, id)
	return &dayBreak, err
}

func (r *BreakRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`
WITH updatedBreak AS (UPDATE break SET is_deleted = true WHERE id = $1),
updatedEvent AS (UPDATE event SET is_deleted = true WHERE break_id = $1)
SELECT TRUE
`, id)
	return err
}

func (r *BreakRepository) Update(dayBreak *entity.Break) error {
	_, err := r.DB.NamedExec(`UPDATE break SET
		day_id = :day_id,
		name = :name,
		start_date = :start_date,
		end_date = :end_date,
		high_bid_team = :high_bid_team,
		giveaway_team = :giveaway_team,
		high_bid_floor = :high_bid_floor
	WHERE id = :id`, dayBreak)

	return err
}

func (r *BreakRepository) GetBreaksByStreamId(dayId int64) ([]*entity.Break, error) {
	breaks := []*entity.Break{}
	err := r.DB.Select(&breaks, `SELECT * FROM break WHERE day_id = $1 AND is_deleted = false`, dayId)
	return breaks, err
}
