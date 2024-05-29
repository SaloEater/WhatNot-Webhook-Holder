package repository_sqlx

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type ChannelRepository struct {
	DB *sqlx.DB
}

func (r *ChannelRepository) GetAll() ([]*entity.Channel, error) {
	channels := []*entity.Channel{}
	err := r.DB.Select(&channels, `SELECT * FROM channel WHERE is_deleted = false`)
	return channels, err
}

func (r *ChannelRepository) Get(id int64) (*entity.Channel, error) {
	var channel entity.Channel
	err := r.DB.Get(&channel, `SELECT * FROM channel where id = $1 AND is_deleted = false`, id)
	return &channel, err
}

func (r *ChannelRepository) GetByStream(streamId int64) (*entity.Channel, error) {
	var channel entity.Channel
	err := r.DB.Get(&channel, `SELECT channel.* FROM channel INNER JOIN public.stream s on channel.id = s.channel_id where s.id = $1 AND channel.is_deleted = false`, streamId)
	return &channel, err
}

func (r *ChannelRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`UPDATE channel SET is_deleted = TRUE WHERE id = $1`, id)
	return err
}

func (r *ChannelRepository) Update(channel *entity.Channel) error {
	_, err := r.DB.NamedExec(`UPDATE channel SET
	  	name = :name,
	  	demo_id = :demo_id
	WHERE id = :id`, channel)

	return err
}

func (r *ChannelRepository) Create(channel *entity.Channel) (int64, error) {
	var id int64
	rows, err := r.DB.NamedQuery(`INSERT INTO channel (
		name, is_deleted
	) VALUES (
	  	:name,
		:is_deleted
	) RETURNING (id)`, channel)
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
