package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type TGChatRepositorier interface {
	CreateOrReEnable(*entity.TGChat) error
	GetByChatID(int64) (*entity.TGChat, error)
	Update(*entity.TGChat) error
	GetAllActive() ([]*entity.TGChat, error)
}
