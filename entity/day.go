package entity

import "time"

type Day struct {
	Id        int64     `db:"id" json:"id"`
	Date      time.Time `db:"date" json:"date"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
}
