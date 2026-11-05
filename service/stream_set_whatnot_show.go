package service

import (
	"time"

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

	oldShows, err := s.WhatnotShowRepositorier.GetByStream(r.StreamId)
	if err != nil {
		return nil, err
	}

	var newShow *entity.WhatnotShow
	if r.WhatnotShowId != nil {
		newShow = &entity.WhatnotShow{
			ChannelId:   stream.ChannelId,
			ShowId:      *r.WhatnotShowId,
			CreatedDate: time.Now().UTC(),
			Status:      entity.WhatnotShowStatusPlanned,
		}
		if r.WhatnotShowName != nil {
			newShow.Title = *r.WhatnotShowName
		}
		if err := s.WhatnotShowRepositorier.Upsert(newShow); err != nil {
			return &SetStreamWhatnotShowResponse{}, err
		}
	}

	// The verification was against the previous show, so it no longer holds. is_ended stays true.
	sameAsBefore := (len(oldShows) == 0 && newShow == nil) || (len(oldShows) == 1 && newShow != nil && oldShows[0].Id == newShow.Id)
	if (stream.Status == entity.StreamStatusVerified || stream.Status == entity.StreamStatusFailedToVerify) && !sameAsBefore {
		stream.Status = entity.StreamStatusEnded
		if err := s.StreamRepositorier.Update(stream); err != nil {
			return &SetStreamWhatnotShowResponse{}, err
		}
		s.StreamCache.Set(cache.IdToKey(r.StreamId), stream)
	}

	var ids []int64
	if newShow != nil {
		ids = []int64{newShow.Id}
	}
	if err := s.WhatnotShowRepositorier.SetStreamShows(r.StreamId, ids); err != nil {
		return &SetStreamWhatnotShowResponse{}, err
	}
	s.StreamShowCache.Delete(cache.IdToKey(r.StreamId))
	if newShow != nil {
		// The upsert may have changed the show's title: other linked streams must reload it.
		s.invalidateShowStreams(newShow.Id)
		s.syncDiscordShowAsync(newShow.Id)
	}
	for _, old := range oldShows {
		if newShow == nil || old.Id != newShow.Id {
			s.syncDiscordShowAsync(old.Id)
		}
	}
	return &SetStreamWhatnotShowResponse{Success: true}, nil
}
