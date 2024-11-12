package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type TrackingRepositorier interface {
	GetValidationData(*entity.ValidationDataBuilder) (*entity.ValidationData, error)
}
