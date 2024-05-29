package service

import (
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetChannelRequest struct {
	Id int64 `json:"id"`
}

func (s *Service) GetChannel(r *GetChannelRequest) (*entity.Channel, error) {
	key := cache.IdToKey(r.Id)

	if !s.ChannelCache.Has(key) {
		channel, err := s.ChannelRepository.Get(r.Id)
		if channel == nil {
			return nil, err
		}
		s.ChannelCache.Set(key, channel)
	}

	cached, found := s.ChannelCache.Get(key)
	if !found {
		return nil, errors.New(fmt.Sprintf("channel for id %d not found", r.Id))
	}

	return cached, nil
}
