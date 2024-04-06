package entity

import "time"

type Break struct {
	Id int64 `json:"id" db:"id"`

	DayId int64  `json:"day_id" db:"day_id"`
	Name  string `json:"name" db:"name"`

	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
}
