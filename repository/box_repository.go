package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type BoxRepositorier interface {
	GetByBundleIDAndItemIndex(int64, int) (*entity.Box, error)
	Update(*entity.Box) error
	BatchCreate([]*entity.Box, int64) error
}
