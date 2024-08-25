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

type GetChannelResponse struct {
	Id                  int64  `json:"id"`
	Name                string `json:"name"`
	DemoId              int64  `json:"demo_id"`
	DefaultHighBidTeam  string `json:"default_high_bid_team"`
	DefaultHighBidFloor int    `json:"default_high_bid_floor"`
}

func (s *Service) GetChannel(r *GetChannelRequest) (*GetChannelResponse, error) {
	s.CacheClear()
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

	response := &GetChannelResponse{
		Id:                  cached.Id,
		Name:                cached.Name,
		DefaultHighBidFloor: cached.DefaultHighBidFloor,
		DefaultHighBidTeam:  cached.DefaultHighBidTeam,
	}

	if cached.DemoId.Valid {
		response.DemoId = cached.DemoId.Int64
	} else {
		response.DemoId = entity.NoId
	}

	return response, nil
}
