package entity

type SeriesTeamTotal struct {
	Team  string `json:"team"  db:"team"`
	Price int64  `json:"price" db:"price"`
}
