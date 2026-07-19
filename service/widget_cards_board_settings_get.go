package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/cache"

type GetWidgetCardsBoardSettingsRequest struct {
	ChannelId int64 `json:"channel_id"`
}

type GetWidgetCardsBoardSettingsResponse struct {
	ChannelId   int64  `json:"channel_id"`
	Orientation string `json:"orientation"`
}

func (s *Service) GetWidgetCardsBoardSettings(r *GetWidgetCardsBoardSettingsRequest) (*GetWidgetCardsBoardSettingsResponse, error) {
	key := cache.IdToKey(r.ChannelId)
	if s.CardsBoardSettingsCache.Has(key) {
		cached, _ := s.CardsBoardSettingsCache.Get(key)
		return &GetWidgetCardsBoardSettingsResponse{ChannelId: cached.ChannelId, Orientation: cached.Orientation}, nil
	}
	w, err := s.WidgetCardsBoardSettingsRepositorier.GetByChannel(r.ChannelId)
	if err != nil {
		return nil, err
	}
	s.CardsBoardSettingsCache.Set(key, w)
	return &GetWidgetCardsBoardSettingsResponse{ChannelId: w.ChannelId, Orientation: w.Orientation}, nil
}
