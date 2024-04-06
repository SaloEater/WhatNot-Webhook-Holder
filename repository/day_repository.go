package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type DayRepository interface {
	GetAll() ([]*entity.Day, error)
	Get(int64) (*entity.Day, error)
	Delete(int64) error
	Update(*entity.Day) error
	Create(*entity.Day) (int64, error)
}
