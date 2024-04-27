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

	response := UpdateDemoResponse{}

	err = s.DemoRepository.Update(demo)

	if err == nil {
		response.Success = true
	}
	s.DemoCache.Set(cache.IdToKey(demo.StreamId), demo)

	return &response, err
}
