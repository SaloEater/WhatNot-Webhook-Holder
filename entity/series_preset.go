package entity

import "time"

type SeriesPreset struct {
	Id        int64     `json:"id"         db:"id"`
	ChannelId int64     `json:"channel_id" db:"channel_id"`
	Name      string    `json:"name"       db:"name"`
	SeriesIds []int64   `json:"series_ids" db:"-"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
