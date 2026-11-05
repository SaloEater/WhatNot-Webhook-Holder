package repository_sqlx

import (
	"database/sql"
	"errors"
	"time"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type WhatnotShowRepository struct {
	DB *sqlx.DB
}

func (r *WhatnotShowRepository) Get(id int64) (*entity.WhatnotShow, error) {
	var v entity.WhatnotShow
	err := r.DB.Unsafe().Get(&v, `SELECT * FROM whatnot_show WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &v, err
}

func (r *WhatnotShowRepository) GetByShowId(channelId int64, showId string) (*entity.WhatnotShow, error) {
	var v entity.WhatnotShow
	err := r.DB.Unsafe().Get(&v, `SELECT * FROM whatnot_show WHERE channel_id = $1 AND show_id = $2`, channelId, showId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &v, err
}

func (r *WhatnotShowRepository) Upsert(show *entity.WhatnotShow) error {
	return r.DB.Get(&show.Id, `INSERT INTO whatnot_show
	(channel_id, show_id, title, description, thumbnail_url, visibility, scheduled_date, created_date, status, updated_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (channel_id, show_id) DO UPDATE SET
	title = CASE WHEN EXCLUDED.title <> '' THEN EXCLUDED.title ELSE whatnot_show.title END,
	description = CASE WHEN EXCLUDED.description <> '' THEN EXCLUDED.description ELSE whatnot_show.description END,
	thumbnail_url = CASE WHEN EXCLUDED.thumbnail_url <> '' THEN EXCLUDED.thumbnail_url ELSE whatnot_show.thumbnail_url END,
	visibility = CASE WHEN EXCLUDED.visibility <> '' THEN EXCLUDED.visibility ELSE whatnot_show.visibility END,
	scheduled_date = COALESCE(EXCLUDED.scheduled_date, whatnot_show.scheduled_date),
	status = CASE
		WHEN whatnot_show.status = 'cancelled' THEN 'cancelled'
		WHEN COALESCE(array_position(ARRAY['planned','live','ended'], EXCLUDED.status::text), 0)
			> COALESCE(array_position(ARRAY['planned','live','ended'], whatnot_show.status::text), 0) THEN EXCLUDED.status
		ELSE whatnot_show.status
	END,
	updated_date = GREATEST(EXCLUDED.updated_date, whatnot_show.updated_date)
RETURNING id`,
		show.ChannelId, show.ShowId, show.Title, show.Description, show.ThumbnailUrl, show.Visibility, show.ScheduledDate, show.CreatedDate, show.Status, show.UpdatedDate)
}

// UpsertOverwrite is Upsert with replace semantics: content fields take the given values, empty ones included.
// An update older than the stored updated_date is skipped in SQL (applied=false, show.Id untouched); a nil
// show.UpdatedDate is always applied and keeps the stored updated_date. Status rules are those of Upsert.
func (r *WhatnotShowRepository) UpsertOverwrite(show *entity.WhatnotShow) (bool, error) {
	var id int64
	err := r.DB.Get(&id, `INSERT INTO whatnot_show
	(channel_id, show_id, title, description, thumbnail_url, visibility, scheduled_date, created_date, status, updated_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (channel_id, show_id) DO UPDATE SET
	title = EXCLUDED.title,
	description = EXCLUDED.description,
	thumbnail_url = EXCLUDED.thumbnail_url,
	visibility = EXCLUDED.visibility,
	scheduled_date = EXCLUDED.scheduled_date,
	status = CASE
		WHEN whatnot_show.status = 'cancelled' THEN 'cancelled'
		WHEN COALESCE(array_position(ARRAY['planned','live','ended'], EXCLUDED.status::text), 0)
			> COALESCE(array_position(ARRAY['planned','live','ended'], whatnot_show.status::text), 0) THEN EXCLUDED.status
		ELSE whatnot_show.status
	END,
	updated_date = COALESCE(EXCLUDED.updated_date, whatnot_show.updated_date)
WHERE whatnot_show.updated_date IS NULL OR EXCLUDED.updated_date IS NULL OR EXCLUDED.updated_date >= whatnot_show.updated_date
RETURNING id`,
		show.ChannelId, show.ShowId, show.Title, show.Description, show.ThumbnailUrl, show.Visibility, show.ScheduledDate, show.CreatedDate, show.Status, show.UpdatedDate)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	show.Id = id
	return true, nil
}

func (r *WhatnotShowRepository) SetStatus(id int64, status string, startedAt *time.Time, endedAt *time.Time) error {
	_, err := r.DB.Exec(`UPDATE whatnot_show SET status = $2, started_at = COALESCE($3, started_at), ended_at = COALESCE($4, ended_at) WHERE id = $1`,
		id, status, startedAt, endedAt)
	return err
}

func (r *WhatnotShowRepository) GetByStream(streamId int64) ([]*entity.WhatnotShow, error) {
	shows := []*entity.WhatnotShow{}
	err := r.DB.Unsafe().Select(&shows, `SELECT w.* FROM whatnot_show w
INNER JOIN stream_whatnot_show l ON l.whatnot_show_id = w.id
WHERE l.stream_id = $1
ORDER BY l.created_at DESC, w.id DESC`, streamId)
	return shows, err
}

func (r *WhatnotShowRepository) GetStreams(whatnotShowId int64) ([]*entity.Stream, error) {
	streams := []*entity.Stream{}
	err := r.DB.Unsafe().Select(&streams, `SELECT s.* FROM stream s
INNER JOIN stream_whatnot_show l ON l.stream_id = s.id
WHERE l.whatnot_show_id = $1 AND s.is_deleted = false
ORDER BY s.id`, whatnotShowId)
	return streams, err
}

func (r *WhatnotShowRepository) GetLatestStream(whatnotShowId int64) (*entity.Stream, error) {
	streams := []*entity.Stream{}
	err := r.DB.Unsafe().Select(&streams, `SELECT s.* FROM stream s
INNER JOIN stream_whatnot_show l ON l.stream_id = s.id
WHERE l.whatnot_show_id = $1 AND s.is_deleted = false
ORDER BY l.created_at DESC, s.id DESC
LIMIT 1`, whatnotShowId)
	if err != nil || len(streams) == 0 {
		return nil, err
	}
	return streams[0], nil
}

func (r *WhatnotShowRepository) SetStreamShows(streamId int64, whatnotShowIds []int64) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM stream_whatnot_show WHERE stream_id = $1`, streamId); err != nil {
		tx.Rollback()
		return err
	}
	for _, id := range whatnotShowIds {
		if _, err := tx.Exec(`INSERT INTO stream_whatnot_show (stream_id, whatnot_show_id) VALUES ($1, $2)`, streamId, id); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
