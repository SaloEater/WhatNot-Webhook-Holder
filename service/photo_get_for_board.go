package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type PhotoGetForBoardRequest struct {
	ChannelId int64 `json:"channel_id"`
	WithSold  bool  `json:"with_sold"`
}

func (s *Service) PhotoGetForBoard(r *PhotoGetForBoardRequest) ([]*entity.Photo, error) {
	return s.PhotoRepositorier.GetUnsoldByChannelActiveStream(r.ChannelId, r.WithSold)
}
