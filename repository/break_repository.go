package repository

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type BreakRepository interface {
	Create(*entity.Break) (int, error)
	Get(int) (*entity.Break, error)
	Delete(int) error
	Update(p *entity.Break) error
}
