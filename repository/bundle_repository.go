package repository

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type BundleRepositorier interface {
	Create(*entity.Bundle) error
	GetByID(int64) (*entity.Bundle, error)
	GetList([]int64) ([]*entity.Bundle, error)
	Update(*entity.Bundle) error
	GetByIDWithLabelUrl(int64) (*entity.BundleWithLabelUrl, error)
}
