package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type CreateWidgetBoardPriceRangeRequest struct {
	ChannelId int64  `json:"channel_id"`
	TierId    string `json:"tier_id"`
	PriceFrom int    `json:"price_from"`
}

type CreateWidgetBoardPriceRangeResponse struct {
	Id      int64 `json:"id"`
	Success bool  `json:"success"`
}

func (s *Service) CreateWidgetBoardPriceRange(r *CreateWidgetBoardPriceRangeRequest) (*CreateWidgetBoardPriceRangeResponse, error) {
	id, err := s.WidgetBoardPriceRangeRepositorier.Create(&entity.WidgetBoardPriceRange{
		ChannelId: r.ChannelId,
		TierId:    r.TierId,
		PriceFrom: r.PriceFrom,
	})
	if err != nil {
		return &CreateWidgetBoardPriceRangeResponse{}, err
	}
	return &CreateWidgetBoardPriceRangeResponse{Id: id, Success: true}, nil
}
