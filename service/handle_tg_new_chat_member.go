package service

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) handleNewChatMember(bot *tgbotapi.BotAPI, update tgbotapi.Update, member tgbotapi.User) {
	if member.IsBot && member.ID == BotID {
		s.setChatForNotifications(bot, update)
	}
}

func (s *Service) setChatForNotifications(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	tgchat := &entity.TGChat{
		ChatID:     update.Message.Chat.ID,
		ChatName:   update.Message.Chat.Title,
		IsDisabled: false,
	}
	err := s.TGChatRepository.CreateOrReEnable(tgchat)
	if err != nil {
		s.sentTGMessage(update.Message.Chat.ID, fmt.Sprintf("Error setting chat for notifications: %s", err.Error()))
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("This chat is registered to recieve notifications about breaks!"))
	_, err = bot.Send(msg)
	if err != nil {
		tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error setting chat for notifications: %s", err.Error()))
		fmt.Println(err)
	}
}
