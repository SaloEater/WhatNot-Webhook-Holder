package entity

type WidgetSeriesBoxesPerBreak struct {
	Id       int64 `json:"id"        db:"id"`
	SeriesId int64 `json:"series_id" db:"series_id"`
	Amount   int   `json:"amount"    db:"amount"`
}
