package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type UpdateWidgetCardsBoardSettingsRequest struct {
	ChannelId   int64  `json:"channel_id"`
	Orientation string `json:"orientation"`
}

type UpdateWidgetCardsBoardSettingsResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateWidgetCardsBoardSettings(r *UpdateWidgetCardsBoardSettingsRequest) (*UpdateWidgetCardsBoardSettingsResponse, error) {
	response := &UpdateWidgetCardsBoardSettingsResponse{}
	err := s.WidgetCardsBoardSettingsRepositorier.Upsert(&entity.WidgetCardsBoardSettings{
		ChannelId:   r.ChannelId,
		Orientation: r.Orientation,
	})
	if err == nil {
		response.Success = true
		s.CardsBoardSettingsCache.Delete(cache.IdToKey(r.ChannelId))
	}
	return response, err
}
