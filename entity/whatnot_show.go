package entity

import "time"

const (
	WhatnotShowStatusPlanned   = "planned"
	WhatnotShowStatusLive      = "live"
	WhatnotShowStatusEnded     = "ended"
	WhatnotShowStatusCancelled = "cancelled"
)

// WhatnotShow is one Whatnot show of a channel. Streams relate to it through stream_whatnot_show.
type WhatnotShow struct {
	Id            int64      `json:"id" db:"id"`
	ChannelId     int64      `json:"channel_id" db:"channel_id"`
	ShowId        string     `json:"show_id" db:"show_id"`
	Title         string     `json:"title" db:"title"`
	Description   string     `json:"description" db:"description"`
	ThumbnailUrl  string     `json:"thumbnail_url" db:"thumbnail_url"`
	Visibility    string     `json:"visibility" db:"visibility"`
	ScheduledDate *time.Time `json:"scheduled_date" db:"scheduled_date"`
	CreatedDate   time.Time  `json:"created_date" db:"created_date"`
	Status        string     `json:"status" db:"status"`
	StartedAt     *time.Time `json:"started_at" db:"started_at"`
	EndedAt       *time.Time `json:"ended_at" db:"ended_at"`
	UpdatedDate   *time.Time `json:"updated_date" db:"updated_date"`
}
