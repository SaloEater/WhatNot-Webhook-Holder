package entity

type Event struct {
	Id int `json:"id"`

	Index      int     `json:"index"`
	Customer   string  `json:"customer"`
	Price      float32 `json:"price"`
	Team       string  `json:"team"`
	IsGiveaway bool    `json:"is_giveaway"`
	Note       string  `json:"note"`
	Quantity   int     `json:"quantity"`
}
