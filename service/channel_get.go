package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetChannelRequest struct {
	Id int64 `json:"id"`
}

func (s *Service) GetChannel(r *GetChannelRequest) (*entity.Channel, error) {
	return s.ChannelRepository.Get(r.Id)
}
