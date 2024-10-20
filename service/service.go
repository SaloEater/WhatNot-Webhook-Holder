package service

import (
	cacheInterface "github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/SaloEater/WhatNot-Webhook-Holder/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	repository.BreakRepository
	repository.StreamRepository
	repository.EventRepository
	repository.DemoRepository
	repository.ChannelRepository
	repository.TGChatRepository
	DemoCache         cacheInterface.Cache[*entity.Demo]
	DemoByStreamCache cacheInterface.Cache[*entity.Demo]
	BreakCache        cacheInterface.Cache[*entity.Break]
	StreamCache       cacheInterface.Cache[*entity.Stream]
	ChannelCache      cacheInterface.Cache[*entity.Channel]
	TelegramBot       *tgbotapi.BotAPI
	StreamShipmenter
}
