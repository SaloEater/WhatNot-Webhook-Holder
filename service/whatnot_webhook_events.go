package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type WhatnotWebhookEventsRequest struct {
	ChannelId int64 `json:"channel_id"`
	Limit     int   `json:"limit"`
}

type WhatnotWebhookEventsResponse struct {
	Items []*entity.WhatnotWebhookEvent `json:"items"`
}

func (s *Service) WhatnotWebhookEvents(r *WhatnotWebhookEventsRequest) (*WhatnotWebhookEventsResponse, error) {
	limit := r.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	items, err := s.ListWebhookEventsByChannel(r.ChannelId, limit)
	if err != nil {
		return nil, err
	}
	return &WhatnotWebhookEventsResponse{Items: items}, nil
}
