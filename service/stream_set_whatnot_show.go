package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/cache"

type SetStreamWhatnotShowRequest struct {
	StreamId      int64   `json:"stream_id"`
	WhatnotShowId *string `json:"whatnot_show_id"`
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
	stream.WhatnotShowId = r.WhatnotShowId
	err = s.StreamRepositorier.Update(stream)
	if err != nil {
		return &SetStreamWhatnotShowResponse{}, err
	}
	s.StreamCache.Set(cache.IdToKey(r.StreamId), stream)
	return &SetStreamWhatnotShowResponse{Success: true}, nil
}
