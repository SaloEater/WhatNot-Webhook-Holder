package entity

type SeriesTeamTotal struct {
	Team       string `json:"team"        db:"team"`
	PriceLeft  int64  `json:"price_left"  db:"price_left"`
	TotalPrice int64  `json:"total_price" db:"total_price"`
}
