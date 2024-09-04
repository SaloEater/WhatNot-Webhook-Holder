package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) handleLeftChatMember(bot *tgbotapi.BotAPI, update tgbotapi.Update, member *tgbotapi.User) {
	if member.IsBot && member.ID == BotID {
		s.disableNotificationsForChat(bot, update)
	}
}

func (s *Service) disableNotificationsForChat(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	tgchat, err := s.TGChatRepository.GetByChatID(update.Message.Chat.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if tgchat == nil {
		fmt.Println(fmt.Sprintf("Chat %d (%s) is not registered in the database!", update.Message.Chat.ID, update.Message.Chat.Title))
		return
	}

	tgchat.IsDisabled = true
	err = s.TGChatRepository.Update(tgchat)
	if err != nil {
		fmt.Println(err)
		return
	}
}
