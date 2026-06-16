package entity

import "time"

type Series struct {
	Id        int64     `json:"id"         db:"id"`
	Name      string    `json:"name"       db:"name"`
	Status    string    `json:"status"     db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
}
