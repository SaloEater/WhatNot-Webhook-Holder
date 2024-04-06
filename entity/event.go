package entity

type Event struct {
	Id int64 `json:"id"`

	BreakId    int64   `db:"break_id"`
	Index      int     `json:"index" db:"index"`
	Customer   string  `json:"customer" db:"customer"`
	Price      float32 `json:"price" db:"price"`
	Team       string  `json:"team" db:"team"`
	IsGiveaway bool    `json:"is_giveaway" db:"is_giveaway"`
	Note       string  `json:"note" db:"note"`
	Quantity   int     `json:"quantity" db:"quantity"`
}
