package service

import (
	"errors"
	"time"
)

const (
	whatnotShowsLookback = 48 * time.Hour
	whatnotShowsPageSize = 20
	whatnotShowsMaxPages = 25
)

type WhatnotShowsListRequest struct {
	ChannelId int64  `json:"channel_id"`
	Status    string `json:"status"`
}

type WhatnotShow struct {
	Id               string     `json:"id"`
	Title            string     `json:"title"`
	Status           string     `json:"status"`
	ScheduledStartAt *time.Time `json:"scheduled_start_at"`
	EndedAt          *time.Time `json:"ended_at"`
}

type WhatnotShowsListResponse struct {
	Shows []WhatnotShow `json:"shows"`
}

func (s *Service) WhatnotShowsList(r *WhatnotShowsListRequest) (*WhatnotShowsListResponse, error) {
	if r.ChannelId == 0 {
		return nil, errors.New("channel_id is required")
	}
	client, err := s.whatnotClientForChannel(r.ChannelId)
	if err != nil {
		return nil, err
	}
	status := r.Status
	if status == "" {
		status = "live"
	}
	// Only recent and upcoming shows are selectable. A live show is always kept, whatever its start date.
	since := time.Now().Add(-whatnotShowsLookback)
	if status == "live" {
		since = time.Time{}
	}
	shows, err := client.ListShowsSince(status, since, whatnotShowsPageSize, whatnotShowsMaxPages)
	if err != nil {
		return nil, err
	}
	resp := &WhatnotShowsListResponse{Shows: make([]WhatnotShow, len(shows))}
	for i, sh := range shows {
		resp.Shows[i] = WhatnotShow{Id: sh.Id, Title: sh.Title, Status: sh.Status, ScheduledStartAt: sh.ScheduledStartAt, EndedAt: sh.EndedAt}
	}
	return resp, nil
}
