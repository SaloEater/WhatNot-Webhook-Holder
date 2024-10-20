package repository_sqlx

import (
	"context"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type StreamShipmentRepository struct {
	DB *sqlx.DB
}

func (s StreamShipmentRepository) GetByStreamID(ctx context.Context, streamID int64) (*entity.StreamShipment, error) {
	shipment := &entity.StreamShipment{}
	err := s.DB.GetContext(ctx, shipment, `
SELECT *
FROM stream_shipment
WHERE stream_id = $1`, streamID)
	return shipment, err
}

func (s StreamShipmentRepository) Update(ctx context.Context, shipment *entity.StreamShipment) error {
	_, err := s.DB.ExecContext(ctx, `
UPDATE stream_shipment
SET task_id = $1, last_task_status = $2
WHERE id = $3`, shipment.TaskID, shipment.LastTaskStatus, shipment.ID)
	return err
}

func (s StreamShipmentRepository) Create(ctx context.Context, shipment *entity.StreamShipment) error {
	_, err := s.DB.NamedExecContext(ctx, `
INSERT INTO stream_shipment (
	stream_id, task_id, last_task_status, stream_shipment_space_id
) VALUES (
	:stream_id, :task_id, :last_task_status, :stream_shipment_space_id
	)`, shipment)
	return err
}
