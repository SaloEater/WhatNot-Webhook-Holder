package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

const (
	StreamShipmentStatusNone = iota
	StreamShipmentStatusStreamInProgress
	StreamShipmentStatusAwaitsLabeling
	StreamShipmentStatusAwaitsPackaging
	StreamShipmentStatusPackagingInProgress
	StreamShipmentStatusAwaitsShipping
	StreamShipmentStatusShipped
)

type StreamShipmenter interface {
	CreateStreamShipment(*entity.Channel, *entity.Stream) error
	MoveStreamShipmentToStatus(int64, int64) error
	DeleteStreamShipment(int64) error
}
