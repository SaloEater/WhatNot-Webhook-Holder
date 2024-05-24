package entity

type Channel struct {
	Id        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	IsDeleted bool   `json:"is_deleted" db:"is_deleted"`
}
