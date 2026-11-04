package repository_sqlx

import (
	"database/sql"
	"errors"
	"time"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type ShowReconciliationRepository struct {
	DB *sqlx.DB
}

func (r *ShowReconciliationRepository) StartReconciliation(streamId int64, showId string, trigger string, staleBefore time.Time) (int64, bool, error) {
	var id int64
	err := r.DB.Get(&id, `INSERT INTO show_reconciliation (stream_id, show_id, trigger, status)
VALUES ($1, $2, $3, 'running')
ON CONFLICT (stream_id, show_id) DO UPDATE SET
	trigger = EXCLUDED.trigger,
	status = 'running',
	started_at = now(),
	finished_at = NULL,
	error = ''
WHERE show_reconciliation.status <> 'running' OR show_reconciliation.started_at < $4
RETURNING id`, streamId, showId, trigger, staleBefore)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return id, true, nil
}

func (r *ShowReconciliationRepository) SetReconciliationManualMap(streamId int64, showId string, manualMapJSON []byte) error {
	_, err := r.DB.Exec(`UPDATE show_reconciliation SET manual_map = $3 WHERE stream_id = $1 AND show_id = $2`, streamId, showId, manualMapJSON)
	return err
}

func (r *ShowReconciliationRepository) SetReconciliationOrders(id int64, ordersJSON []byte) error {
	_, err := r.DB.Exec(`UPDATE show_reconciliation SET orders = $2 WHERE id = $1`, id, ordersJSON)
	return err
}

func (r *ShowReconciliationRepository) FinishReconciliation(id int64, status string, errMsg string, reportJSON []byte) error {
	_, err := r.DB.Exec(`UPDATE show_reconciliation SET status = $2, error = $3, report = $4, finished_at = now() WHERE id = $1`, id, status, errMsg, reportJSON)
	return err
}

func (r *ShowReconciliationRepository) GetReconciliation(id int64) (*entity.ShowReconciliation, error) {
	var v entity.ShowReconciliation
	err := r.DB.Unsafe().Get(&v, `SELECT * FROM show_reconciliation WHERE id = $1`, id)
	return &v, err
}

func (r *ShowReconciliationRepository) GetReconciliationByStream(streamId int64, showId string) (*entity.ShowReconciliation, error) {
	var v entity.ShowReconciliation
	err := r.DB.Unsafe().Get(&v, `SELECT * FROM show_reconciliation WHERE stream_id = $1 AND show_id = $2`, streamId, showId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &v, err
}
