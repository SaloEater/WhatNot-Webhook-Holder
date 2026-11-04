package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type StreamReconcileRequest struct {
	StreamId int64 `json:"stream_id"`
	DryRun   bool  `json:"dry_run"` // report only: nothing is applied
}

// StreamReconcile runs a manual reconciliation synchronously. A failed run is returned as the row
// (status "failed", error set) rather than as an HTTP error.
func (s *Service) StreamReconcile(r *StreamReconcileRequest) (*entity.ShowReconciliation, error) {
	trigger := "manual"
	if r.DryRun {
		trigger = "report"
	}
	row, err := s.RunShowReconciliation(r.StreamId, trigger, r.DryRun)
	if row != nil {
		return row, nil
	}
	return nil, err
}
