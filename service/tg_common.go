package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) sentTGMessage(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := s.TelegramBot.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}