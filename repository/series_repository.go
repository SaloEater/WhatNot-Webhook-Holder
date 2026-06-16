package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type SeriesRepositorier interface {
	Create(name string) (int64, error)
	Get(id int64) (*entity.Series, error)
	GetList() ([]*entity.Series, error)
	Update(id int64, name string) error
	Close(id int64) error
	Delete(id int64) error
	CountOpen() (int, error)
}
