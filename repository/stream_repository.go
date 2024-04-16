package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type StreamRepository interface {
	GetAll() ([]*entity.Stream, error)
	Get(int64) (*entity.Stream, error)
	Delete(int64) error
	Update(*entity.Stream) error
	Create(*entity.Stream) (int64, error)
}