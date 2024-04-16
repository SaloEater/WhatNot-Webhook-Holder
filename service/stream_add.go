package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"time"
)

type AddStreamRequest struct {
	Name string `json:"name"`
}

type AddStreamResponse struct {
	Id int64 `json:"id"`
}

func (s *Service) AddStream(r *AddStreamRequest) (*AddStreamResponse, error) {
	id, err := s.StreamRepository.Create(&entity.Stream{
		Name:      r.Name,
		CreatedAt: time.Now().UTC(),
		IsDeleted: false,
	})

	if err != nil {
		return nil, err
	}

	return &AddStreamResponse{Id: id}, nil
}
