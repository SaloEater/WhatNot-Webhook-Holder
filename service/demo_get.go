package service

import (
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetDemoRequest struct {
	Id int64 `json:"id"`
}

type GetDemoResponse struct {
	Id                int64  `db:"id"`
	StreamId          int64  `json:"stream_id"`
	BreakId           int64  `json:"break_id"`
	HighlightUsername string `json:"highlight_username"`
}

func (s *Service) GetDemo(r *GetDemoRequest) (*GetDemoResponse, error) {
	key := cache.IdToKey(r.Id)
	var err error
	var breakId int64

	if !s.DemoCache.Has(key) {
		demo, err := s.DemoRepositorier.Get(r.Id)
		if demo == nil {
			return nil, err
		}
		s.DemoCache.Set(key, demo)
		streamKey := cache.IdToKey(demo.StreamId)
		s.DemoByStreamCache.Set(streamKey, demo)
	}

	cached, found := s.DemoCache.Get(key)
	if !found {
		return nil, errors.New(fmt.Sprintf("demo for id %d not found", r.Id))
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
