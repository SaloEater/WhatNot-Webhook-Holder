package service

import (
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
)

type GetChannelRequest struct {
	Id int64 `json:"id"`
}

type GetChannelResponse struct {
	Id                  int64  `json:"id"`
	Name                string `json:"name"`
	ActiveStreamId      *int64 `json:"active_stream_id"`
	DefaultHighBidTeam  string `json:"default_high_bid_team"`
	DefaultHighBidFloor int    `json:"default_high_bid_floor"`
}

func (s *Service) GetChannel(r *GetChannelRequest) (*GetChannelResponse, error) {
	key := cache.IdToKey(r.Id)

	if !s.ChannelCache.Has(key) {
		channel, err := s.ChannelRepositorier.Get(r.Id)
		if channel == nil {
			return nil, err
		}
		s.ChannelCache.Set(key, channel)
	}

	cached, found := s.ChannelCache.Get(key)
	if !found {
		return nil, errors.New(fmt.Sprintf("channel for id %d not found", r.Id))
	}

	return &GetChannelResponse{
		Id:                  cached.Id,
		Name:                cached.Name,
		ActiveStreamId:      cached.ActiveStreamId,
		DefaultHighBidFloor: cached.DefaultHighBidFloor,
		DefaultHighBidTeam:  cached.DefaultHighBidTeam,
	}, nil
}
