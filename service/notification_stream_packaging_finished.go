package service

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type NotificationStreamPackagingFinished struct {
	StreamID int64 `json:"stream_id"`
}

func (s *Service) NotificationStreamPackagingFinished(r *NotificationStreamPackagingFinished) error {
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

	return nil
}

func getStreamCreatedTIme(stream *entity.StreamEnriched) string {
	return stream.CreatedAt.Format("2006-01-02")
}
