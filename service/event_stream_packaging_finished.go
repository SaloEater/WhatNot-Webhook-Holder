package service

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/pkg/errors"
)

type EventStreamPackagingFinished struct {
	StreamID int64 `json:"stream_id"`
}

func (s *Service) EventStreamPackagingFinished(r *EventStreamPackagingFinished) error {
	tgchats, err := s.TGChatRepository.GetAllActive()
	if err != nil {
		return err
	}

	stream, err := s.StreamRepository.GetEnriched(r.StreamID)
	if err != nil {
		return err
	}

	for _, tgchat := range tgchats {
		s.sendTGMessage(tgchat.ChatID, fmt.Sprintf("Packaging finished for stream <b>%s</b> from <b>%s</b> for (%s)!", stream.Name, stream.ChannelName, getStreamCreatedTIme(stream)))
	}

	go func() {
		err = s.MoveStreamShipmentToStatus(stream.Id, StreamShipmentStatusAwaitsShipping)
		if err != nil {
			fmt.Println(errors.WithMessage(err, fmt.Sprintf("move %d (%s) stream to awaits packaging", stream.Id, stream.ChannelName)))
		}
	}()

	return nil
}

func getStreamCreatedTIme(stream *entity.StreamEnriched) string {
	return stream.CreatedAt.Format("2006-01-02")
}
