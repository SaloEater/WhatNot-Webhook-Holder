package service

import (
	"database/sql"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"time"
)

type AddStreamRequest struct {
	Name string `json:"name"`
}

type AddStreamResponse struct {
	GetStreamsStream
}

func (s *Service) AddStream(r *AddStreamRequest) (*AddStreamResponse, error) {
	stream := entity.Stream{
		Name:      r.Name,
		CreatedAt: time.Now().UTC(),
		IsDeleted: false,
	}
	id, err := s.StreamRepository.Create(&stream)
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

	return &AddStreamResponse{GetStreamsStream: GetStreamsStream{
		Id:        stream.Id,
		Name:      stream.Name,
		CreatedAt: stream.CreatedAt.UnixMilli(),
	}}, nil
}
