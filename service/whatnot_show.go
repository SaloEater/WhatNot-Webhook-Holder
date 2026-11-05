package service

import (
	"time"

	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

// streamShow returns the stream's most recently linked show, nil when it has none.
func (s *Service) streamShow(streamId int64) (*entity.WhatnotShow, error) {
	key := cache.IdToKey(streamId)
	if cached, found := s.StreamShowCache.Get(key); found {
		return cached, nil
	}
	shows, err := s.WhatnotShowRepositorier.GetByStream(streamId)
	if err != nil {
		return nil, err
	}
	var show *entity.WhatnotShow
	if len(shows) > 0 {
		show = shows[0]
	}
	s.StreamShowCache.Set(key, show)
	return show, nil
}

// fillStreamShow sets the JSON-only whatnot_show_id / whatnot_show_name fields from the link.
func (s *Service) fillStreamShow(stream *entity.Stream) error {
	show, err := s.streamShow(stream.Id)
	if err != nil {
		return err
	}
	stream.WhatnotShowId, stream.WhatnotShowName = nil, nil
	if show != nil {
		id := show.ShowId
		stream.WhatnotShowId = &id
		if show.Title != "" {
			name := show.Title
			stream.WhatnotShowName = &name
		}
	}
	return nil
}

// invalidateShowStreams drops the cached show of every stream linked to the show; they reload lazily.
func (s *Service) invalidateShowStreams(whatnotShowId int64) {
	streams, err := s.WhatnotShowRepositorier.GetStreams(whatnotShowId)
	if err != nil {
		s.StreamShowCache.Clear()
		return
	}
	for _, st := range streams {
		s.StreamShowCache.Delete(cache.IdToKey(st.Id))
	}
}

// ensureWhatnotShow returns the show row, creating it when missing. A created row takes title and
// schedule from Whatnot (GetShow); when that fails a bare row with an empty title is created.
func (s *Service) ensureWhatnotShow(channelId int64, showId string) (*entity.WhatnotShow, error) {
	show, err := s.WhatnotShowRepositorier.GetByShowId(channelId, showId)
	if err != nil || show != nil {
		return show, err
	}
	show = &entity.WhatnotShow{ChannelId: channelId, ShowId: showId, Status: entity.WhatnotShowStatusPlanned, CreatedDate: time.Now().UTC()}
	if client, cerr := s.whatnotClientForChannel(channelId); cerr == nil {
		if wn, gerr := client.GetShow(showId); gerr == nil {
			wn.Id = showId
			show = whatnotShowFromAPI(channelId, wn)
			// The started/ended/cancelled handlers set the status (and its timestamps) themselves.
			show.Status = entity.WhatnotShowStatusPlanned
		}
	}
	if err := s.WhatnotShowRepositorier.Upsert(show); err != nil {
		return nil, err
	}
	return s.WhatnotShowRepositorier.Get(show.Id)
}

// whatnotShowStatusOrder ranks the non-final statuses; cancelled is outside the order.
func whatnotShowStatusOrder(status string) int {
	switch status {
	case entity.WhatnotShowStatusLive:
		return 1
	case entity.WhatnotShowStatusEnded:
		return 2
	}
	return 0
}
