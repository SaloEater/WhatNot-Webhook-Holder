package service

import (
	"encoding/json"
	"fmt"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

func (s *Service) streamForWebhookShow(ev *entity.WhatnotWebhookEvent, data json.RawMessage) (*entity.Stream, error) {
	var d struct {
		ShowId string `json:"show_id"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	if d.ShowId == "" {
		return nil, fmt.Errorf("no show_id in event data")
	}
	stream, err := s.StreamRepositorier.GetByWhatnotShowId(ev.ChannelId, d.ShowId)
	if err != nil {
		return nil, err
	}
	if stream == nil {
		return nil, fmt.Errorf("no stream linked to show %s", d.ShowId)
	}
	return stream, nil
}

func (s *Service) onWhatnotShowStarted(ev *entity.WhatnotWebhookEvent, data json.RawMessage) error {
	stream, err := s.streamForWebhookShow(ev, data)
	if err != nil {
		return err
	}
	switch stream.Status {
	case entity.StreamStatusLive, entity.StreamStatusEnded, entity.StreamStatusVerified:
		return nil
	}
	return s.SetStreamStatus(stream, entity.StreamStatusLive)
}

func (s *Service) onWhatnotShowEnded(ev *entity.WhatnotWebhookEvent, data json.RawMessage) error {
	stream, err := s.streamForWebhookShow(ev, data)
	if err != nil {
		return err
	}
	if err := s.EventStreamEnded(&EventStreamEndedRequest{StreamID: stream.Id}); err != nil {
		fmt.Println(fmt.Sprintf("whatnot webhook: stream %d: event stream ended: %s", stream.Id, err.Error()))
	}
	_, err = s.RunShowReconciliation(stream.Id, "webhook")
	return err
}
