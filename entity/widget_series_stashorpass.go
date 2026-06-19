package entity

type WidgetSeriesStashorpass struct {
	Id        int64 `json:"id" db:"id"`
	ChannelId int64 `json:"channel_id" db:"channel_id"`
	Price     int   `json:"price" db:"price"`
}
