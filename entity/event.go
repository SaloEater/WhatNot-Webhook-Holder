package entity

type Event struct {
	Id string `json:"id"`

	Customer   string  `json:"customer"`
	Price      float32 `json:"price"`
	Team       string  `json:"team"`
	IsGiveaway bool    `json:"is_giveaway"`
	Note       string  `json:"note"`
	Quantity   int     `json:"quantity"`
}
