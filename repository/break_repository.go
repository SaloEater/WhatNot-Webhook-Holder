package repository

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type BreakRepository interface {
	Create(*entity.Break) (int64, error)
	Get(int64) (*entity.Break, error)
	Delete(int64) error
	Update(p *entity.Break) error
	GetBreaksByDay(int64) ([]*entity.Break, error)
}
