package repository_sqlx

import (
	"fmt"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type PhotoRepository struct {
	DB *sqlx.DB
}

func (r *PhotoRepository) Create(p *entity.Photo) (int64, error) {
	var id int64
	rows, err := r.DB.NamedQuery(`INSERT INTO photo (
		series_id, name, team, url
	) VALUES (
		:series_id,
		:name,
		:team,
		:url
	) RETURNING (id)`, p)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("no rows returned after INSERT")
	}
	return id, nil
}

func (r *PhotoRepository) GetBySeriesId(seriesId int64) ([]*entity.Photo, error) {
	photos := []*entity.Photo{}
	err := r.DB.Select(&photos, `SELECT * FROM photo WHERE series_id = $1`, seriesId)
	return photos, err
}

func (r *PhotoRepository) GetUnsoldByChannelActiveStream(channelId int64) ([]*entity.Photo, error) {
	photos := []*entity.Photo{}
	err := r.DB.Select(&photos, `
SELECT p.*
FROM photo p
JOIN series s ON s.id = p.series_id
JOIN break b ON b.series_id = s.id
JOIN stream st ON st.id = b.day_id
JOIN channel c ON c.active_stream_id = st.id
WHERE c.id = $1
--   AND s.status = 'closed'
  AND p.is_sold = false
  AND p.is_deleted = false
  AND s.is_deleted = false
  AND b.is_deleted = false
`, channelId)
	return photos, err
}

func (r *PhotoRepository) Update(id int64, name, team string) error {
	_, err := r.DB.Exec(`UPDATE photo SET name = $1, team = $2 WHERE id = $3`, name, team, id)
	return err
}

func (r *PhotoRepository) MarkSold(id int64, sold bool) error {
	_, err := r.DB.Exec(`UPDATE photo SET is_sold = $1 WHERE id = $2`, sold, id)
	return err
}

func (r *PhotoRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`UPDATE photo SET is_deleted = true WHERE id = $1`, id)
	return err
}

func (r *PhotoRepository) Restore(id int64) error {
	_, err := r.DB.Exec(`UPDATE photo SET is_deleted = false WHERE id = $1`, id)
	return err
}
