package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type SetStreamWhatnotShowRequest struct {
	StreamId        int64   `json:"stream_id"`
	WhatnotShowId   *string `json:"whatnot_show_id"`
	WhatnotShowName *string `json:"whatnot_show_name"`
}

type SetStreamWhatnotShowResponse struct {
	Success bool `json:"success"`
}

func (s *Service) SetStreamWhatnotShow(r *SetStreamWhatnotShowRequest) (*SetStreamWhatnotShowResponse, error) {
	stream, err := s.StreamRepositorier.Get(r.StreamId)
	if err != nil {
		return nil, err
	}
	if r.WhatnotShowId != nil && *r.WhatnotShowId == "" {
		r.WhatnotShowId = nil
	}
	if r.WhatnotShowId == nil {
		r.WhatnotShowName = nil
	}
	// The verification was against the previous show, so it no longer holds. is_ended stays true.
	if stream.Status == entity.StreamStatusVerified && !sameShow(stream.WhatnotShowId, r.WhatnotShowId) {
		stream.Status = entity.StreamStatusEnded
	}
	stream.WhatnotShowId = r.WhatnotShowId
	stream.WhatnotShowName = r.WhatnotShowName
	err = s.StreamRepositorier.Update(stream)
	if err != nil {
		return &SetStreamWhatnotShowResponse{}, err
	}
	s.StreamCache.Set(cache.IdToKey(r.StreamId), stream)
	return &SetStreamWhatnotShowResponse{Success: true}, nil
}

func sameShow(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
