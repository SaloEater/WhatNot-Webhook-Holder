package repository_sqlx

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	DB *sqlx.DB
}

func (r *EventRepository) GetAllByBreak(breakId int64) ([]*entity.Event, error) {
	var events []*entity.Event
	err := r.DB.Select(&events, `SELECT * FROM event WHERE break_id = $1`, breakId)
	return events, err
}

func (r *EventRepository) Get(id int64) (*entity.Event, error) {
	var event entity.Event
	err := r.DB.Get(&event, `SELECT * FROM event where id = $1`, id)
	return &event, err
}

func (r *EventRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM event WHERE id = $1`, id)
	return err
}

func (r *EventRepository) Update(event *entity.Event) error {
	_, err := r.DB.NamedExec(`UPDATE event SET
		break_id = :break_id,
		index = :index,
		customer = :customer,
		price = :price,
		team = :team,
		is_giveaway = :is_giveaway,
		note = :note,
		quantity = :quantity
	WHERE id = :id`, event)

	return err
}

func (r *EventRepository) Create(event *entity.Event) (int64, error) {
	var id int64
	rows, err := r.DB.NamedQuery(`INSERT INTO event (
		break_id,
		index,
		customer,
		price,
		team,
		is_giveaway,
		note,
		quantity
	) VALUES (
	          :break_id,
	          (SELECT COALESCE(MAX(index), 0) + 1 FROM event WHERE break_id = :break_id),
	          :customer,
	          :price,
	          :team,
	          :is_giveaway,
	          :note,
	          :quantity
	) RETURNING (id)`, event)
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

func (r *EventRepository) GetAllChildren(eventId int64) ([]*entity.Event, error) {
	events := []*entity.Event{}
	err := r.DB.Select(&events, `
		SELECT * FROM event WHERE break_id IN (
		    SELECT break_id FROM event WHERE id = $1
		)
	`, eventId)
	return events, err
}

func (r *EventRepository) UpdateAll(events []*entity.Event) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, event := range events {
		_, err = r.DB.NamedExec(`UPDATE event SET
			break_id = :break_id,
			index = :index,
			customer = :customer,
			price = :price,
			team = :team,
			is_giveaway = :is_giveaway,
			note = :note,
			quantity = :quantity
		WHERE id = :id`, event)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
