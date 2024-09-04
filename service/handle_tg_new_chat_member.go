package service

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) handleNewChatMember(update tgbotapi.Update, member tgbotapi.User) {
	if member.IsBot && member.ID == BotID {
		s.setChatForNotifications(update)
	}
}

func (s *Service) setChatForNotifications(update tgbotapi.Update) {
	tgchat := &entity.TGChat{
		ChatID:     update.Message.Chat.ID,
		ChatName:   update.Message.Chat.Title,
		IsDisabled: false,
	}
	err := s.TGChatRepository.CreateOrReEnable(tgchat)
	if err != nil {
		s.sendTGMessage(update.Message.Chat.ID, fmt.Sprintf("Error setting chat for notifications: %s", err.Error()))
		return
	}

	s.sendTGMessage(update.Message.Chat.ID, fmt.Sprintf("This chat is registered to recieve notifications about breaks!"))
}
