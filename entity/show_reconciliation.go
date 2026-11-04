package entity

import (
	"encoding/json"
	"time"
)

const (
	ReconciliationStatusRunning     = "running"
	ReconciliationStatusApplied     = "applied"
	ReconciliationStatusNeedsManual = "needs_manual"
	ReconciliationStatusFailed      = "failed"
)

type ShowReconciliation struct {
	Id         int64           `json:"id" db:"id"`
	StreamId   int64           `json:"stream_id" db:"stream_id"`
	ShowId     string          `json:"show_id" db:"show_id"`
	Trigger    string          `json:"trigger" db:"trigger"`
	Status     string          `json:"status" db:"status"`
	StartedAt  time.Time       `json:"started_at" db:"started_at"`
	FinishedAt *time.Time      `json:"finished_at" db:"finished_at"`
	Error      string          `json:"error" db:"error"`
	Orders     json.RawMessage `json:"orders" db:"orders"`
	Report     json.RawMessage `json:"report" db:"report"`
	ManualMap  json.RawMessage `json:"manual_map" db:"manual_map"`
}
