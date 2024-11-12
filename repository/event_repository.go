package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type EventRepositorier interface {
	GetAllByBreak(int64) ([]*entity.Event, error)
	Get(int64) (*entity.Event, error)
	Delete(int64) error
	Update(*entity.Event) error
	Create(*entity.Event) (int64, error)
	GetAllChildren(int64) ([]*entity.Event, error)
	UpdateAll([]*entity.Event) error
	Move(int64, int) error
	GetAll([]int64) ([]*entity.Event, error)
	GetAvailableByChannelIDAndTeam(int64, string) (*entity.Event, error)
	GetLastIndex(int64) (int, error)
}
