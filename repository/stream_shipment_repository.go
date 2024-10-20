package repository

import (
	"context"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type StreamShipmentRepository interface {
	GetByStreamID(context.Context, int64) (*entity.StreamShipment, error)
	Update(context.Context, *entity.StreamShipment) error
	Create(context.Context, *entity.StreamShipment) error
}
