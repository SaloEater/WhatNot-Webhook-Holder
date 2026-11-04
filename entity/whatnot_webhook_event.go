package entity

import (
	"encoding/json"
	"time"
)

const (
	WebhookEventStatusReceived  = "received"
	WebhookEventStatusProcessed = "processed"
	WebhookEventStatusIgnored   = "ignored"
	WebhookEventStatusFailed    = "failed"
)

type WhatnotWebhookEvent struct {
	Id          int64           `json:"id" db:"id"`
	WebhookId   string          `json:"webhook_id" db:"webhook_id"`
	ChannelId   int64           `json:"channel_id" db:"channel_id"`
	Type        string          `json:"type" db:"type"`
	SellerId    string          `json:"seller_id" db:"seller_id"`
	EventTime   *time.Time      `json:"event_time" db:"event_time"`
	ReceivedAt  time.Time       `json:"received_at" db:"received_at"`
	Payload     json.RawMessage `json:"payload,omitempty" db:"payload"`
	Status      string          `json:"status" db:"status"`
	Error       string          `json:"error" db:"error"`
	ProcessedAt *time.Time      `json:"processed_at" db:"processed_at"`
}
