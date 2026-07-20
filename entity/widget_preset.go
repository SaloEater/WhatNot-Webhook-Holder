package entity

type WidgetPreset struct {
	Id        int64  `json:"id"         db:"id"`
	ChannelId int64  `json:"channel_id" db:"channel_id"`
	Name      string `json:"name"       db:"name"`
	Value     string `json:"value"      db:"value"`
}
