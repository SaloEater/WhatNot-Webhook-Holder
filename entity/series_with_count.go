package entity

type SeriesWithCount struct {
	*Series
	UnsoldCount int64 `json:"unsold_count"`
	SoldCount   int64 `json:"sold_count"`
}
