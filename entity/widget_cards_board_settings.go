package entity

type WidgetCardsBoardSettings struct {
	Id          int64  `json:"id"          db:"id"`
	ChannelId   int64  `json:"channel_id"  db:"channel_id"`
	Orientation string `json:"orientation" db:"orientation"`
}
