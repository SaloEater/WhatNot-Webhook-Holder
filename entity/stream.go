package entity

import "time"

type Stream struct {
	Id        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	ChannelId int64     `json:"channel_id" db:"channel_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
}
