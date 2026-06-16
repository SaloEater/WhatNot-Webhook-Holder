package entity

type SeriesTeamPrice struct {
	Id       int64   `json:"id"        db:"id"`
	SeriesId int64   `json:"series_id" db:"series_id"`
	Team     string  `json:"team"      db:"team"`
	Price    float64 `json:"price"     db:"price"`
}
