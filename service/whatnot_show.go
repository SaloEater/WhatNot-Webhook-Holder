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

// stampShowStatusTime fills started_at / ended_at when an upsert moved the status to live or ended without the
// started/ended event (the full-show payloads carry the status but no timestamps). The event time is the best guess.
func (s *Service) stampShowStatusTime(whatnotShowId int64, at time.Time) error {
	show, err := s.WhatnotShowRepositorier.Get(whatnotShowId)
	if err != nil || show == nil {
		return err
	}
	switch {
	case show.Status == entity.WhatnotShowStatusLive && show.StartedAt == nil:
		return s.WhatnotShowRepositorier.SetStatus(show.Id, show.Status, &at, nil)
	case show.Status == entity.WhatnotShowStatusEnded && show.EndedAt == nil:
		return s.WhatnotShowRepositorier.SetStatus(show.Id, show.Status, nil, &at)
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
// A known row is refreshed from Whatnot (overwrite semantics, updated_date guard); the fetched status is
// replaced by the row's current one so the refresh never moves status. A failed lookup leaves it unchanged.
// When the refresh was applied the stream caches are dropped and a Discord sync is started.
func (s *Service) ensureWhatnotShow(channelId int64, showId string) (*entity.WhatnotShow, error) {
	show, err := s.WhatnotShowRepositorier.GetByShowId(channelId, showId)
	if err != nil {
		return nil, err
	}
	if show != nil {
		client, cerr := s.whatnotClientForChannel(channelId)
		if cerr != nil {
			return show, nil
		}
		wn, gerr := client.GetShow(showId)
		if gerr != nil || wn == nil {
			return show, nil
		}
		wn.Id = showId
		fresh := whatnotShowFromAPI(channelId, wn)
		fresh.Status = show.Status
		applied, uerr := s.WhatnotShowRepositorier.UpsertOverwrite(fresh)
		if uerr != nil {
			return nil, uerr
		}
		if !applied {
			return show, nil
		}
		s.invalidateShowStreams(show.Id)
		s.syncDiscordShowAsync(show.Id)
		return s.WhatnotShowRepositorier.Get(show.Id)
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
