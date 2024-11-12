package repository

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type BundleLabelRepositorier interface {
	Create(labels *entity.BundleLabels) error
}
