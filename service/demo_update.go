package service

import (
	"database/sql"
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type UpdateDemoRequest struct {
	Id                int64  `json:"id"`
	BreakId           int64  `json:"break_id"`
	HighlightUsername string `json:"highlight_username"`
}

type UpdateDemoResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateDemo(r *UpdateDemoRequest) (*UpdateDemoResponse, error) {
	demo, err := s.DemoRepository.Get(r.Id)
	if err != nil {
		return nil, err
	}

	if r.BreakId == entity.NoBreakId {
		demo.BreakId = sql.NullInt64{Valid: false}
	} else {
		demo.BreakId = sql.NullInt64{Int64: r.BreakId, Valid: true}
	}
	demo.HighlightUsername = r.HighlightUsername

	err = s.DemoRepository.Update(demo)

	if err != nil {
		return &UpdateDemoResponse{Success: false}, err
	}
	s.DemoCache.Set(cache.IdToKey(demo.StreamId), demo)

	err = s.updateChannelDemo(demo.Id)
	if err != nil {
		return &UpdateDemoResponse{Success: false}, err
	}

	return &UpdateDemoResponse{Success: true}, err
}

func (s *Service) updateChannelDemo(streamId int64) error {
	channel, err := s.ChannelRepository.GetByStream(streamId)
	if err != nil {
		return err
	}

	channel.DemoId = sql.NullInt64{Valid: true, Int64: streamId}
	err = s.ChannelRepository.Update(channel)
	if err != nil {
		return err
	}

	key := cache.IdToKey(channel.Id)
	s.ChannelCache.Set(key, channel)
	return nil
}
