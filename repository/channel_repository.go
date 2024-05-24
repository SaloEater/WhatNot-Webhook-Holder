package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type ChannelRepository interface {
	GetAll() ([]*entity.Channel, error)
	Get(int64) (*entity.Channel, error)
	Delete(int64) error
	Update(*entity.Channel) error
	Create(*entity.Channel) (int64, error)
	GetByStream(int64) (*entity.Channel, error)
}
