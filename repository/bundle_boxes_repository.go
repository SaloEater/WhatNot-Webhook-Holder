package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type BundleBoxesRepositorier interface {
	GetAllByBundle(int64) ([]*entity.BundleBoxes, error)
	Create(*entity.BundleBoxes) error
	GetByID(int64) (*entity.BundleBoxes, error)
	Update(*entity.BundleBoxes) error
	Delete(int64) error
	GetByBundleIDAndBoxTypeID(int64, int64) (*entity.BundleBoxes, error)
}
