package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type DayRepository interface {
	GetAll() ([]*entity.Day, error)
	Get(int) (*entity.Day, error)
	Delete(int) error
	Update(*entity.Day) error
	Create(*entity.Day) (int, error)
}
