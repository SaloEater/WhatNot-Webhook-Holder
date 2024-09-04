package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NotificationStreamEndedRequest struct {
	StreamID int64 `json:"stream_id"`
}

func (s *Service) NotificationStreamEnded(r *NotificationStreamEndedRequest) error {
	tgchats, err := s.TGChatRepository.GetAll()
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
		`Stream %s from %s has ended!
Breaks: %d
Earned: %d$
Unique customers: %d`,
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

		msg := tgbotapi.NewMessage(tgchat.ChatID, statsMessage)
		_, err := s.TelegramBot.Send(msg)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
