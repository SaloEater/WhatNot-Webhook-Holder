package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type EventRepository interface {
	GetAllByBreak(int) ([]*entity.Event, error)
	Get(int) (*entity.Event, error)
	Delete(int) error
	Update(*entity.Event) error
	Create(*entity.Event) (int, error)
	GetAllChildren(int) ([]*entity.Event, error)
	UpdateAll([]*entity.Event) error
}
