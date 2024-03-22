package entity

type Break struct {
	Name string `json:"name"`

	Events []Event `json:"events"`

	StartDate int64 `json:"start_date"`
	EndDate   int64 `json:"end_date"`
}
