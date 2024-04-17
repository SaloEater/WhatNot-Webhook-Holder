package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

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
	demo, err := s.DemoRepository.GetByStream(r.StreamId)
	if demo == nil {
		return nil, err
	}

	var breakId int64
	if demo.BreakId.Valid {
		breakId = demo.BreakId.Int64
	} else {
		breakId = entity.NoBreakId
	}

	return &GetDemoResponse{
		Id:                demo.Id,
		StreamId:          demo.StreamId,
		BreakId:           breakId,
		HighlightUsername: demo.HighlightUsername,
	}, err
}
