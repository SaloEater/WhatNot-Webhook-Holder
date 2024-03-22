package entity

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type Day struct {
	Date   Date     `json:"date"`
	Breaks []string `json:"breaks"`
}
