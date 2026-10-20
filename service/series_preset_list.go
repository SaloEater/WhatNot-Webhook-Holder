package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type SeriesPresetListRequest struct {
	ChannelId int64 `json:"channel_id"`
}

type SeriesPresetListResponse struct {
	Presets []*entity.SeriesPreset `json:"presets"`
}

func (s *Service) SeriesPresetList(r *SeriesPresetListRequest) (*SeriesPresetListResponse, error) {
	presets, err := s.SeriesPresetRepositorier.ListByChannel(r.ChannelId)
	if err != nil {
		return nil, err
	}
	return &SeriesPresetListResponse{Presets: presets}, nil
}
