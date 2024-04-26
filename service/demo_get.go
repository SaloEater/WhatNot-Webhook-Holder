package service

import (
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetDemoRequest struct {
	StreamId int64 `json:"stream_id"`
}

type GetDemoResponse struct {
	Id                int64  `db:"id" json:"id"`
	StreamId          int64  `json:"stream_id" db:"stream_id"`
	BreakId           int64  `json:"break_id" db:"break_id"`
	HighlightUsername string `json:"highlight_username" db:"highlight_username"`
}

func (s *Service) GetDemo(r *GetDemoRequest) (*GetDemoResponse, error) {
	key := cache.IdToKey(r.StreamId)
	var err error
	var breakId int64

	if !s.DemoCache.Has(key) {
		demo, err := s.DemoRepository.GetByStream(r.StreamId)
		if demo == nil {
			return nil, err
		}
		s.DemoCache.Set(key, demo)
	}

	cached, found := s.DemoCache.Get(key)
	if !found {
		return nil, errors.New(fmt.Sprintf("demo for stream %d not found", r.StreamId))
	}

	if cached.BreakId.Valid {
		breakId = cached.BreakId.Int64
	} else {
		breakId = entity.NoBreakId
	}

	return &GetDemoResponse{
		Id:                cached.Id,
		StreamId:          cached.StreamId,
		BreakId:           breakId,
		HighlightUsername: cached.HighlightUsername,
	}, err
}
