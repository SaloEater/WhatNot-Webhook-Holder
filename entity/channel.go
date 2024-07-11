package entity

import "database/sql"

type Channel struct {
	Id                  int64         `json:"id" db:"id"`
	Name                string        `json:"name" db:"name"`
	IsDeleted           bool          `json:"is_deleted" db:"is_deleted"`
	DemoId              sql.NullInt64 `json:"demo_id" db:"demo_id"`
	DefaultHighBidTeam  string        `json:"default_high_bid_team" db:"default_high_bid_team"`
	DefaultHighBidFloor int           `json:"default_high_bid_floor" db:"default_high_bid_floor"`
}
