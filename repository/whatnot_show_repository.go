package repository

import (
	"time"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type WhatnotShowRepositorier interface {
	// Get and GetByShowId return nil, nil when there is no such show.
	Get(id int64) (*entity.WhatnotShow, error)
	GetByShowId(channelId int64, showId string) (*entity.WhatnotShow, error)
	// Upsert is idempotent per (channel_id, show_id): empty fields never overwrite stored values and the
	// status never moves backwards (planned < live < ended; cancelled is final). Sets show.Id.
	Upsert(show *entity.WhatnotShow) error
	// UpsertOverwrite is Upsert with replace semantics for title, description, thumbnail, visibility and scheduled
	// date (empty values included). An event older than the stored updated_date is skipped atomically: applied is
	// false then and no error is returned. A nil UpdatedDate is always applied. Sets show.Id when applied.
	UpsertOverwrite(show *entity.WhatnotShow) (applied bool, err error)
	// SetStatus sets the status; non-nil startedAt/endedAt are stored, nil ones keep the stored value.
	SetStatus(id int64, status string, startedAt *time.Time, endedAt *time.Time) error
	// GetByStream returns the stream's shows, most recently linked first.
	GetByStream(streamId int64) ([]*entity.WhatnotShow, error)
	// GetStreams returns the non-deleted streams linked to the show.
	GetStreams(whatnotShowId int64) ([]*entity.Stream, error)
	// GetLatestStream returns the non-deleted stream most recently linked to the show (latest link created_at, then
	// highest stream id); nil, nil when there is none.
	GetLatestStream(whatnotShowId int64) (*entity.Stream, error)
	// SetStreamShows replaces the stream's link rows.
	SetStreamShows(streamId int64, whatnotShowIds []int64) error
}
