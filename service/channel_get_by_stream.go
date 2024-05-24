package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetChannelByStreamRequest struct {
	StreamId int64 `json:"stream_id"`
}

func (s *Service) GetChannelByStream(r *GetChannelByStreamRequest) (*entity.Channel, error) {
	return s.ChannelRepository.GetByStream(r.StreamId)
}
