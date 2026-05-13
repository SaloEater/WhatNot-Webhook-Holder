package service

import (
	cacheInterface "github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/SaloEater/WhatNot-Webhook-Holder/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	repository.BreakRepositorier
	repository.StreamRepositorier
	repository.EventRepositorier
	repository.ChannelRepositorier
	repository.TGChatRepositorier
	repository.BundleRepositorier
	repository.BoxTypeRepositorier
	repository.LocationRepositorier
	repository.BundleBoxesRepositorier
	repository.BoxRepositorier
	repository.TrackingRepositorier
	repository.BundleLabelRepositorier
	BreakCache   cacheInterface.Cache[*entity.Break]
	StreamCache  cacheInterface.Cache[*entity.Stream]
	ChannelCache cacheInterface.Cache[*entity.Channel]
	TelegramBot  *tgbotapi.BotAPI
	StreamShipmenter
	DigitalOceaner
}
