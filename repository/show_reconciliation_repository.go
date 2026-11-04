package repository

import (
	"time"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type ShowReconciliationRepositorier interface {
	// StartReconciliation creates the row of this stream + show, or resets it for a re-run (the manual map
	// survives). started is false when a run that began after staleBefore is still in progress.
	StartReconciliation(streamId int64, showId string, trigger string, staleBefore time.Time) (id int64, started bool, err error)
	SetReconciliationOrders(id int64, ordersJSON []byte) error
	SetReconciliationManualMap(streamId int64, showId string, manualMapJSON []byte) error
	FinishReconciliation(id int64, status string, errMsg string, reportJSON []byte) error
	GetReconciliation(id int64) (*entity.ShowReconciliation, error)
	// GetReconciliationByStream returns nil, nil when none.
	GetReconciliationByStream(streamId int64, showId string) (*entity.ShowReconciliation, error)
}
