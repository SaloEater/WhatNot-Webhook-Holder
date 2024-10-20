package service

import (
	"database/sql"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/pkg/errors"
	"time"
)

type AddStreamRequest struct {
	Name      string `json:"name"`
	ChannelId int64  `json:"channel_id"`
}

type AddStreamResponse struct {
	GetStreamResponse
}

func (s *Service) AddStream(r *AddStreamRequest) (*AddStreamResponse, error) {
	stream := &entity.Stream{
		Name:      r.Name,
		CreatedAt: time.Now().UTC(),
		IsDeleted: false,
		ChannelId: r.ChannelId,
	}
	id, err := s.StreamRepository.Create(stream)
	if err != nil {
		return nil, err
	}

	stream.Id = id

	_, err = s.DemoRepository.Create(&entity.Demo{
		StreamId:          stream.Id,
		BreakId:           sql.NullInt64{Valid: false},
		HighlightUsername: "",
	})
	if err != nil {
		return nil, err
	}

	go func() {
		channel, err := s.ChannelRepository.Get(r.ChannelId)
		if err != nil {
			fmt.Println(errors.WithMessage(err, "get channel by id"))
		}
		err = s.CreateStreamShipment(channel, stream)
		if err != nil {
			fmt.Println(errors.WithMessage(err, "create stream shipment"))
		}
	}()

	return &AddStreamResponse{GetStreamResponse: GetStreamResponse{
		Id:        stream.Id,
		Name:      stream.Name,
		CreatedAt: stream.CreatedAt.UnixMilli(),
	}}, nil
}
