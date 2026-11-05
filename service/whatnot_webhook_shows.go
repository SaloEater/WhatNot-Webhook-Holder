package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/SaloEater/WhatNot-Webhook-Holder/whatnot"
)

func webhookShowId(data json.RawMessage) (string, error) {
	var d struct {
		ShowId string `json:"show_id"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return "", err
	}
	if d.ShowId == "" {
		return "", fmt.Errorf("no show_id in event data")
	}
	return d.ShowId, nil
}

func webhookEventTime(ev *entity.WhatnotWebhookEvent) time.Time {
	if ev.EventTime != nil {
		return ev.EventTime.UTC()
	}
	return time.Now().UTC()
}

func parseWebhookTime(v string) *time.Time {
	t, err := time.Parse(time.RFC3339Nano, v)
	if err != nil {
		return nil
	}
	return &t
}

// whatnotShowFromAPI maps the Whatnot show shape (REST get-a-show and the data of shows.created/updated) to a row.
// Unknown statuses map to planned, which the forward-only status rule never lets regress a stored status.
func whatnotShowFromAPI(channelId int64, wn *whatnot.Show) *entity.WhatnotShow {
	show := &entity.WhatnotShow{
		ChannelId:     channelId,
		ShowId:        wn.Id,
		Title:         wn.Title,
		Description:   wn.Description,
		Visibility:    wn.Visibility,
		ScheduledDate: wn.ScheduledStartAt,
		CreatedDate:   time.Now().UTC(),
		UpdatedDate:   parseWebhookTime(wn.UpdatedAt),
		Status:        entity.WhatnotShowStatusPlanned,
	}
	if t := parseWebhookTime(wn.CreatedAt); t != nil {
		show.CreatedDate = *t
	}
	if wn.Thumbnail != nil {
		show.ThumbnailUrl = wn.Thumbnail.Url
	}
	switch wn.Status {
	case "live":
		show.Status = entity.WhatnotShowStatusLive
	case "ended":
		show.Status = entity.WhatnotShowStatusEnded
	case "cancelled":
		show.Status = entity.WhatnotShowStatusCancelled
	}
	return show
}

func decodeWhatnotShowEvent(ev *entity.WhatnotWebhookEvent, data json.RawMessage) (*entity.WhatnotShow, error) {
	var wn whatnot.Show
	if err := json.Unmarshal(data, &wn); err != nil {
		return nil, err
	}
	if wn.Id == "" {
		return nil, fmt.Errorf("no id in event data")
	}
	return whatnotShowFromAPI(ev.ChannelId, &wn), nil
}

func (s *Service) onWhatnotShowCreated(ev *entity.WhatnotWebhookEvent, data json.RawMessage) error {
	show, err := decodeWhatnotShowEvent(ev, data)
	if err != nil {
		return err
	}
	if err := s.WhatnotShowRepositorier.Upsert(show); err != nil {
		return err
	}
	s.invalidateShowStreams(show.Id)
	s.syncDiscordShowAsync(show.Id)
	return nil
}

// onWhatnotShowUpdated overwrites the content fields; a stale event (older updated_at) is ignored without a sync.
func (s *Service) onWhatnotShowUpdated(ev *entity.WhatnotWebhookEvent, data json.RawMessage) error {
	show, err := decodeWhatnotShowEvent(ev, data)
	if err != nil {
		return err
	}
	applied, err := s.WhatnotShowRepositorier.UpsertOverwrite(show)
	if err != nil || !applied {
		return err
	}
	s.invalidateShowStreams(show.Id)
	s.syncDiscordShowAsync(show.Id)
	return nil
}

func (s *Service) onWhatnotShowStarted(ev *entity.WhatnotWebhookEvent, data json.RawMessage) error {
	showId, err := webhookShowId(data)
	if err != nil {
		return err
	}
	show, err := s.ensureWhatnotShow(ev.ChannelId, showId)
	if err != nil {
		return err
	}
	if show.Status == entity.WhatnotShowStatusPlanned {
		startedAt := webhookEventTime(ev)
		if err := s.WhatnotShowRepositorier.SetStatus(show.Id, entity.WhatnotShowStatusLive, &startedAt, nil); err != nil {
			return err
		}
		s.invalidateShowStreams(show.Id)
		s.syncDiscordShowAsync(show.Id)
	}

	streams, err := s.WhatnotShowRepositorier.GetStreams(show.Id)
	if err != nil {
		return err
	}
	for _, stream := range streams {
		switch stream.Status {
		case entity.StreamStatusLive, entity.StreamStatusEnded, entity.StreamStatusFailedToVerify, entity.StreamStatusVerified:
			continue
		}
		if err := s.SetStreamStatus(stream, entity.StreamStatusLive); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) onWhatnotShowEnded(ev *entity.WhatnotWebhookEvent, data json.RawMessage) error {
	showId, err := webhookShowId(data)
	if err != nil {
		return err
	}
	show, err := s.ensureWhatnotShow(ev.ChannelId, showId)
	if err != nil {
		return err
	}
	if show.Status != entity.WhatnotShowStatusCancelled && show.Status != entity.WhatnotShowStatusEnded {
		endedAt := webhookEventTime(ev)
		if err := s.WhatnotShowRepositorier.SetStatus(show.Id, entity.WhatnotShowStatusEnded, nil, &endedAt); err != nil {
			return err
		}
		s.invalidateShowStreams(show.Id)
		s.syncDiscordShowAsync(show.Id)
	}

	streams, err := s.WhatnotShowRepositorier.GetStreams(show.Id)
	if err != nil {
		return err
	}
	var firstErr error
	for _, stream := range streams {
		if err := s.EventStreamEnded(&EventStreamEndedRequest{StreamID: stream.Id}); err != nil {
			fmt.Println(fmt.Sprintf("whatnot webhook: stream %d: event stream ended: %s", stream.Id, err.Error()))
		}
		if _, err := s.RunShowReconciliation(stream.Id, "webhook", false); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// onWhatnotShowCancelled does nothing for a show we have no row for: the payload only carries the id.
func (s *Service) onWhatnotShowCancelled(ev *entity.WhatnotWebhookEvent, data json.RawMessage) error {
	showId, err := webhookShowId(data)
	if err != nil {
		return err
	}
	show, err := s.WhatnotShowRepositorier.GetByShowId(ev.ChannelId, showId)
	if err != nil || show == nil {
		return err
	}
	if err := s.WhatnotShowRepositorier.SetStatus(show.Id, entity.WhatnotShowStatusCancelled, nil, nil); err != nil {
		return err
	}
	s.invalidateShowStreams(show.Id)
	s.syncDiscordShowAsync(show.Id)
	return nil
}
