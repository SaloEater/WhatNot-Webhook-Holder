package service

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
)

type NotificationStreamEndedRequest struct {
	StreamID int64 `json:"stream_id"`
}

func (s *Service) NotificationStreamEnded(r *NotificationStreamEndedRequest) error {
	tgchats, err := s.TGChatRepository.GetAllActive()
	if err != nil {
		return err
	}

	streamStats, err := s.StreamRepository.GetStats(r.StreamID)
	if err != nil {
		return err
	}
	if streamStats == nil {
		return fmt.Errorf("stream not found")
	}

	statsMessage := fmt.Sprintf(
		`Stream <b>%s</b> from <b>%s</b> has ended!
Breaks: <b>%d</b>
Earned: <b>%d$</b>
Unique customers: <b>%d</b>`,
		streamStats.Name,
		streamStats.ChannelName,
		streamStats.BreaksAmount,
		streamStats.SoldFor,
		streamStats.CustomersAmount,
	)

	if len(streamStats.BigCustomers) > 0 {
		statsMessage += "\nBig customers: "
		for _, customer := range streamStats.BigCustomers {
			statsMessage += fmt.Sprintf("%s, ", customer)
		}
	}

	if len(streamStats.LuckyGoblins) > 0 {
		statsMessage += "\nLucky giveaway goblins: "
		for _, goblin := range streamStats.LuckyGoblins {
			statsMessage += fmt.Sprintf("%s, ", goblin)
		}
	}

	for _, tgchat := range tgchats {
		if tgchat.IsDisabled {
			continue
		}

		s.sendTGMessage(tgchat.ChatID, statsMessage)
	}

	stream, err := s.StreamRepository.Get(r.StreamID)
	if err != nil {
		return err
	}

	stream.IsEnded = true
	err = s.StreamRepository.Update(stream)
	key := cache.IdToKey(stream.Id)
	s.StreamCache.Set(key, stream)

	return nil
}
