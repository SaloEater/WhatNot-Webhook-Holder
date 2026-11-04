package service

import (
	"fmt"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/pkg/errors"
)

type EventStreamEndedRequest struct {
	StreamID int64 `json:"stream_id"`
}

func (s *Service) EventStreamEnded(r *EventStreamEndedRequest) error {
	stream, err := s.StreamRepositorier.GetEnriched(r.StreamID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("get enriched stream %d", r.StreamID))
	}

	// Idempotent: a webhook retry or manual-then-webhook must not repeat the shipment move / TG message.
	if stream.IsEnded {
		return nil
	}

	err = s.SetStreamStatus(stream.Stream, entity.StreamStatusEnded)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("set stream %d status to ended", r.StreamID))
	}

	go func() {
		err := s.MoveStreamShipmentToStatus(stream.Id, StreamShipmentStatusAwaitsLabeling)
		if err != nil {
			fmt.Println(errors.WithMessage(err, fmt.Sprintf("move %d (%s) stream to awaits labeling", stream.Id, stream.ChannelName)))
		}
	}()

	tgchats, err := s.TGChatRepositorier.GetAllActive()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("get all active tgchats for stream %d", r.StreamID))
	}

	streamStats, err := s.StreamRepositorier.GetStats(r.StreamID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("get stream %d stats", r.StreamID))
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

	return nil
}
