package repository

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type BreakRepositorier interface {
	Create(*entity.Break) (int64, error)
	Get(int64) (*entity.Break, error)
	Delete(int64) error
	Update(p *entity.Break) error
	GetBreaksByStreamId(int64) ([]*entity.Break, error)
}
