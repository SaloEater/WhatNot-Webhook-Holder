package entity

type WidgetBoardPriceRange struct {
	Id        int64  `json:"id"         db:"id"`
	ChannelId int64  `json:"channel_id" db:"channel_id"`
	TierId    string `json:"tier_id"    db:"tier_id"`
	PriceFrom int    `json:"price_from" db:"price_from"`
}
