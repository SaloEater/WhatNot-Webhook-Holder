package entity

type Channel struct {
	Id                  int64  `json:"id" db:"id"`
	Name                string `json:"name" db:"name"`
	IsDeleted           bool   `json:"is_deleted" db:"is_deleted"`
	ActiveStreamId      *int64 `json:"active_stream_id" db:"active_stream_id"`
	DefaultHighBidTeam  string `json:"default_high_bid_team" db:"default_high_bid_team"`
	DefaultHighBidFloor int    `json:"default_high_bid_floor" db:"default_high_bid_floor"`
}
