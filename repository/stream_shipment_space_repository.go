package repository

import (
	"context"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type StreamShipmentSpaceRepository interface {
	GetByChannelID(context.Context, int64) (*entity.StreamShipmentSpace, error)
}
