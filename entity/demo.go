package entity

import "database/sql"

const NoBreakId = 0

type Demo struct {
	Id                int64         `db:"id" json:"id"`
	StreamId          int64         `json:"stream_id" db:"stream_id"`
	BreakId           sql.NullInt64 `json:"break_id" db:"break_id"`
	HighlightUsername string        `json:"highlight_username" db:"highlight_username"`
}
