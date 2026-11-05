package entity

import "time"

// WhatnotShowId/WhatnotShowName are not columns: they are filled from the stream_whatnot_show link (Service.fillStreamShow).
type Stream struct {
	Id              int64     `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	ChannelId       int64     `json:"channel_id" db:"channel_id"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	IsDeleted       bool      `json:"is_deleted" db:"is_deleted"`
	IsEnded         bool      `json:"is_ended" db:"is_ended"`
	ActiveBreakId   *int64    `json:"active_break_id" db:"active_break_id"`
	WhatnotShowId   *string   `json:"whatnot_show_id" db:"-"`
	WhatnotShowName *string   `json:"whatnot_show_name" db:"-"`
	Status          string    `json:"status" db:"status"`
}

const (
	StreamStatusPlanned  = "planned"
	StreamStatusLive     = "live"
	StreamStatusEnded    = "ended"
	StreamStatusVerified = "verified"
	// StreamStatusFailedToVerify sits between ended and verified: the last real reconciliation was needs_manual or failed.
	StreamStatusFailedToVerify = "failed_to_verify"
)

type StreamStatistic struct {
	Name            string   `db:"name"`
	BreaksAmount    int      `db:"breaks_amount"`
	SoldFor         int      `db:"sold_for"`
	CustomersAmount int      `db:"customers_amount"`
	ChannelName     string   `db:"channel_name"`
	BigCustomers    []string `db:"big_customers"`
	LuckyGoblins    []string `db:"lucky_goblins"`
}

type StreamEnriched struct {
	*Stream
	ChannelName string `db:"channel_name"`
}
