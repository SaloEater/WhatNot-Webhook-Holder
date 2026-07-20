package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type ListWidgetPresetsRequest struct {
	ChannelId int64 `json:"channel_id"`
}

type ListWidgetPresetsResponse struct {
	Presets []*entity.WidgetPreset `json:"presets"`
}

func (s *Service) ListWidgetPresets(r *ListWidgetPresetsRequest) (*ListWidgetPresetsResponse, error) {
	presets, err := s.WidgetPresetRepositorier.ListByChannel(r.ChannelId)
	if err != nil {
		return nil, err
	}
	return &ListWidgetPresetsResponse{Presets: presets}, nil
}
