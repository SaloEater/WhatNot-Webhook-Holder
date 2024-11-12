package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type BoxTypeRepositorier interface {
	GetAll() ([]*entity.BoxType, error)
	GetByID(int64) (*entity.BoxType, error)
	Create(*entity.BoxType) error
	Update(*entity.BoxType) error
	GetByIDs([]int64) ([]*entity.BoxType, error)
}
