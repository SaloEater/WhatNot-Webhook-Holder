package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type StreamReconcileRequest struct {
	StreamId int64 `json:"stream_id"`
}

// StreamReconcile runs a manual reconciliation synchronously. A failed run is returned as the row
// (status "failed", error set) rather than as an HTTP error.
func (s *Service) StreamReconcile(r *StreamReconcileRequest) (*entity.ShowReconciliation, error) {
	row, err := s.RunShowReconciliation(r.StreamId, "manual")
	if row != nil {
		return row, nil
	}
	return nil, err
}
