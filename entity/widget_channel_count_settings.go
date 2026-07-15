package entity

type WidgetChannelCountSettings struct {
	Id             int64 `json:"id"              db:"id"`
	ChannelId      int64 `json:"channel_id"      db:"channel_id"`
	ShowPercentage bool  `json:"show_percentage" db:"show_percentage"`
}
