package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type ListWidgetBoardPriceRangesRequest struct {
	ChannelId int64 `json:"channel_id"`
}

type ListWidgetBoardPriceRangesResponse struct {
	Ranges []*entity.WidgetBoardPriceRange `json:"ranges"`
}

var defaultPriceRanges = []*entity.WidgetBoardPriceRange{
	{TierId: "best", PriceFrom: 700},
	{TierId: "good", PriceFrom: 450},
	{TierId: "mid", PriceFrom: 50},
}

func (s *Service) ListWidgetBoardPriceRanges(r *ListWidgetBoardPriceRangesRequest) (*ListWidgetBoardPriceRangesResponse, error) {
	ranges, err := s.WidgetBoardPriceRangeRepositorier.ListByChannel(r.ChannelId)
	if err != nil {
		return nil, err
	}
	if len(ranges) == 0 {
		defaults := make([]*entity.WidgetBoardPriceRange, len(defaultPriceRanges))
		for i, d := range defaultPriceRanges {
			defaults[i] = &entity.WidgetBoardPriceRange{ChannelId: r.ChannelId, TierId: d.TierId, PriceFrom: d.PriceFrom}
		}
		return &ListWidgetBoardPriceRangesResponse{Ranges: defaults}, nil
	}
	return &ListWidgetBoardPriceRangesResponse{Ranges: ranges}, nil
}
