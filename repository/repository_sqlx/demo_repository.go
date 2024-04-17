package repository_sqlx

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type DemoRepository struct {
	DB *sqlx.DB
}

func (r *DemoRepository) Get(id int64) (*entity.Demo, error) {
	var demo entity.Demo
	err := r.DB.Get(&demo, `SELECT * FROM demo WHERE id = $1`, id)
	return &demo, err
}

func (r *DemoRepository) GetByStream(streamId int64) (*entity.Demo, error) {
	var demo entity.Demo
	err := r.DB.Get(&demo, `SELECT * FROM demo WHERE stream_id = $1 LIMIT 1`, streamId)
	return &demo, err
}

func (r *DemoRepository) Update(demo *entity.Demo) error {
	_, err := r.DB.NamedExec(`UPDATE demo SET
		break_id = :break_id, highlight_username = :highlight_username
	WHERE id = :id`, demo)

	return err
}

func (r *DemoRepository) Create(demo *entity.Demo) (int64, error) {
	var id int64
	rows, err := r.DB.NamedQuery(`INSERT INTO demo (
		break_id, highlight_username, stream_id
	) VALUES (
		:break_id,
		:highlight_username,
		:stream_id
	) RETURNING (id)`, demo)
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
