package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type DemoRepository interface {
	Get(int64) (*entity.Demo, error)
	GetByStream(int64) (*entity.Demo, error)
	Update(*entity.Demo) error
	Create(*entity.Demo) (int64, error)
}
