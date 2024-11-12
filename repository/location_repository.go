package repository

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type LocationRepositorier interface {
	GetAll() ([]*entity.Location, error)
	GetByID(int64) (*entity.Location, error)
}
