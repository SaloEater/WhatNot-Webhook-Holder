package service

import (
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetDemoByStreamRequest struct {
	StreamId int64 `json:"stream_id"`
}

func (s *Service) GetDemoByStream(r *GetDemoByStreamRequest) (*GetDemoResponse, error) {
	key := cache.IdToKey(r.StreamId)
	var err error
	var breakId int64

	if !s.DemoByStreamCache.Has(key) {
		demo, err := s.DemoRepositorier.GetByStream(r.StreamId)
		if demo == nil {
			return nil, err
		}
		s.DemoByStreamCache.Set(key, demo)
		idKey := cache.IdToKey(demo.Id)
		s.DemoCache.Set(idKey, demo)
	}

	cached, found := s.DemoByStreamCache.Get(key)
	if !found {
		return nil, errors.New(fmt.Sprintf("demo for stream %d not found", r.StreamId))
	}

	if cached.BreakId.Valid {
		breakId = cached.BreakId.Int64
	} else {
		breakId = entity.NoId
	}

	return &GetDemoResponse{
		Id:                cached.Id,
		StreamId:          cached.StreamId,
		BreakId:           breakId,
		HighlightUsername: cached.HighlightUsername,
	}, err
}
