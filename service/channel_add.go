package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type AddChannelRequest struct {
	Name string `json:"name"`
}

func (s *Service) AddChannel(r *AddChannelRequest) (*GetChannelsChannel, error) {
	stream := entity.Channel{
		Name:      r.Name,
		IsDeleted: false,
	}
	id, err := s.ChannelRepositorier.Create(&stream)
	if err != nil {
		return nil, err
	}

	stream.Id = id

	return &GetChannelsChannel{
		Id:   stream.Id,
		Name: stream.Name,
	}, nil
}
