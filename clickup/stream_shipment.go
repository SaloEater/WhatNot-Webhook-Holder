package clickup

import (
	"context"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/SaloEater/WhatNot-Webhook-Holder/repository"
	"github.com/SaloEater/WhatNot-Webhook-Holder/repository/repository_sqlx"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"github.com/jmoiron/sqlx"
	"github.com/raksul/go-clickup/clickup"
)

const StreamDateFieldID = "0d14d7e7-d4c0-40a8-8b13-da1a77ca1fab"

type StreamShipment struct {
	client *clickup.Client
	repository.StreamShipmentRepository
	repository.StreamShipmentSpaceRepository
}

func Init(apiKey string, db *sqlx.DB) *StreamShipment {
	client := clickup.NewClient(nil, apiKey)
	return &StreamShipment{
		client:                        client,
		StreamShipmentRepository:      &repository_sqlx.StreamShipmentRepository{DB: db},
		StreamShipmentSpaceRepository: &repository_sqlx.StreamShipmentSpaceRepository{DB: db},
	}
}

func (s *StreamShipment) CreateStreamShipment(channel *entity.Channel, stream *entity.Stream) error {
	ctx := context.Background()
	shipmentSpace, err := s.StreamShipmentSpaceRepository.GetByChannelID(ctx, channel.Id)
	if err != nil {
		return err
	}

	task, _, err := s.client.Tasks.CreateTask(ctx, fmt.Sprintf("%d", shipmentSpace.ListID), &clickup.TaskRequest{
		Name:   fmt.Sprintf("%s #%d", stream.Name, stream.Id),
		Status: getStatus(service.StreamShipmentStatusStreamInProgress),
		CustomFields: []clickup.CustomFieldInTaskRequest{
			{ID: StreamDateFieldID, Value: fmt.Sprintf("%d", stream.CreatedAt.UnixMilli())},
		},
	})
	if err != nil {
		return err
	}
	shipment := &entity.StreamShipment{
		StreamID:              stream.Id,
		TaskID:                task.ID,
		LastTaskStatus:        task.Status.Status,
		StreamShipmentSpaceID: shipmentSpace.ID,
	}
	return s.StreamShipmentRepository.Create(ctx, shipment)
}

func getStatus(status int64) string {
	switch status {
	case service.StreamShipmentStatusStreamInProgress:
		return "stream in progress"
	case service.StreamShipmentStatusAwaitsLabeling:
		return "awaits labeling"
	case service.StreamShipmentStatusAwaitsShipping:
		return "awaits shipping"
	default:
		return "stream in progress"
	}
}

func (s *StreamShipment) MoveStreamShipmentToStatus(streamId, status int64) error {
	ctx := context.Background()
	shipment, err := s.StreamShipmentRepository.GetByStreamID(ctx, streamId)
	if err != nil {
		return err
	}
	task, _, err := s.client.Tasks.UpdateTask(ctx, shipment.TaskID, &clickup.GetTaskOptions{}, &clickup.TaskUpdateRequest{
		Status: getStatus(status),
	})
	if err != nil {
		return err
	}
	shipment.LastTaskStatus = task.Status.Status
	return s.StreamShipmentRepository.Update(ctx, shipment)
}

func (s *StreamShipment) DeleteStreamShipment(streamId int64) error {
	ctx := context.Background()
	shipment, err := s.GetByStreamID(ctx, streamId)
	if err != nil {
		return err
	}
	_, err = s.client.Tasks.DeleteTask(ctx, shipment.TaskID, &clickup.GetTaskOptions{})
	if err != nil {
		return err
	}
	shipment.TaskID = ""
	return s.StreamShipmentRepository.Update(ctx, shipment)
}
