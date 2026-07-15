package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type UpdateWidgetChannelCountSettingsRequest struct {
	ChannelId      int64 `json:"channel_id"`
	ShowPercentage bool  `json:"show_percentage"`
}

type UpdateWidgetChannelCountSettingsResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateWidgetChannelCountSettings(r *UpdateWidgetChannelCountSettingsRequest) (*UpdateWidgetChannelCountSettingsResponse, error) {
	response := &UpdateWidgetChannelCountSettingsResponse{}
	err := s.WidgetChannelCountSettingsRepositorier.Upsert(&entity.WidgetChannelCountSettings{
		ChannelId:      r.ChannelId,
		ShowPercentage: r.ShowPercentage,
	})
	if err == nil {
		response.Success = true
	}
	return response, err
}
