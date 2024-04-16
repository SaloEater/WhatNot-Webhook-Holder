package entity

import "time"

type Stream struct {
	Id        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
}
