package repository_sqlx

import (
	"context"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type StreamShipmentSpaceRepository struct {
	DB *sqlx.DB
}

func (r *StreamShipmentSpaceRepository) GetByChannelID(ctx context.Context, channelID int64) (*entity.StreamShipmentSpace, error) {
	var shipment entity.StreamShipmentSpace
	err := r.DB.Get(&shipment, `SELECT * FROM stream_shipment_space WHERE channel_id = $1`, channelID)
	return &shipment, err
}
