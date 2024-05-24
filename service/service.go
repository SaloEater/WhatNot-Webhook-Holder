package service

import (
	cacheInterface "github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/SaloEater/WhatNot-Webhook-Holder/repository"
)

type Service struct {
	repository.BreakRepository
	repository.StreamRepository
	repository.EventRepository
	repository.DemoRepository
	repository.ChannelRepository
	DemoCache  cacheInterface.Cache[*entity.Demo]
	BreakCache cacheInterface.Cache[*entity.Break]
}
