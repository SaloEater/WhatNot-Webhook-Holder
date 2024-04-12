package repository_sqlx

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	DB *sqlx.DB
}

func (r *EventRepository) GetAll(ids []int64) ([]*entity.Event, error) {
	query, args, err := sqlx.Named(`SELECT * FROM event WHERE id IN (:ids)`, map[string]interface{}{"ids": ids})
	if err != nil {
		return nil, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = r.DB.Rebind(query)

	events := make([]*entity.Event, len(ids))
	err = r.DB.Select(&events, query, args...)
	return events, err
}

func (r *EventRepository) GetAllByBreak(breakId int64) ([]*entity.Event, error) {
	var events []*entity.Event
	err := r.DB.Select(&events, `SELECT * FROM event WHERE break_id = $1 AND is_deleted = false`, breakId)
	return events, err
}

func (r *EventRepository) Get(id int64) (*entity.Event, error) {
	var event entity.Event
	err := r.DB.Get(&event, `SELECT * FROM event where id = $1 AND is_deleted = false`, id)
	return &event, err
}

func (r *EventRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`UPDATE event SET is_deleted = true WHERE id = $1`, id)
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
		quantity,
        is_deleted
	) VALUES (
		:break_id,
		(SELECT COALESCE(MAX(index), 0) + 1 FROM event WHERE break_id = :break_id),
		:customer,
		:price,
		:team,
		:is_giveaway,
		:note,
		:quantity,
		:is_deleted
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
		) AND is_deleted = false
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
		WHERE id = :id;`, event)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *EventRepository) Move(id int64, newIndex int) error {
	tx, err := r.DB.Beginx()
	defer func() {
		errRb := tx.Rollback()
		if errRb != nil {
			fmt.Println("trying to cancel committed transaction")
		}
	}()
	query := `
WITH oldRecord AS (
    SELECT index as oldIndex, break_id as breakId FROM event WHERE id = :id
) UPDATE event
SET
    index = CASE
            WHEN :new_index > oldRecord.oldIndex THEN
                CASE
                    WHEN index <= :new_index AND index > oldRecord.oldIndex THEN index - 1
                    ELSE index
                END
            WHEN :new_index < oldRecord.oldIndex THEN
                CASE
                    WHEN index >= :new_index AND index < oldRecord.oldIndex THEN index + 1
                    ELSE index
                END
            ELSE index
        END
FROM oldRecord
WHERE break_id = oldRecord.breakId;
`
	args := map[string]interface{}{
		"id":        id,
		"new_index": newIndex,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`UPDATE event SET index = :new_index WHERE id = :id;`, args)
	if err != nil {
		return err
	}

	return tx.Commit()
}
