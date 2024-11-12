package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type BundleRepository struct {
	DB *sqlx.DB
}

func (r *BundleRepository) Create(bundle *entity.Bundle) error {
	var id int64
	err := r.DB.QueryRow(
		`INSERT INTO bundle (location_id, status, is_deleted, name) VALUES ($1, $2, $3, $4) RETURNING id`,
		bundle.LocationID, bundle.Status, bundle.IsDeleted, bundle.Name,
	).Scan(&id)
	if err == nil {
		bundle.ID = id
	}
	return err
}

func (r *BundleRepository) GetByID(id int64) (*entity.Bundle, error) {
	var bundle entity.Bundle
	err := r.DB.Get(&bundle, `SELECT * FROM bundle WHERE id = $1 and is_deleted = false`, id)
	return &bundle, err
}

func (r *BundleRepository) GetList(locationIds []int64) ([]*entity.Bundle, error) {
	var bundles []*entity.Bundle
	var err error

	if len(locationIds) == 0 {
		err = r.DB.Select(&bundles, `SELECT * FROM bundle WHERE is_deleted = false`)
	} else {
		query, args, err := sqlx.In(`SELECT * FROM bundle WHERE location_id IN (?) AND is_deleted = false`, locationIds)
		if err != nil {
			return nil, err
		}
		query = r.DB.Rebind(query)
		err = r.DB.Select(&bundles, query, args...)
	}

	return bundles, err
}

func (r *BundleRepository) Update(bundle *entity.Bundle) error {
	_, err := r.DB.Exec(
		`UPDATE bundle SET location_id = $1, status = $2, name = $3, locked_at = $4 WHERE id = $5`,
		bundle.LocationID, bundle.Status, bundle.Name, bundle.LockedAt, bundle.ID,
	)
	return err
}

func (r *BundleRepository) GetByIDWithLabelUrl(id int64) (*entity.BundleWithLabelUrl, error) {
	var bundle entity.BundleWithLabelUrl
	err := r.DB.Get(&bundle, `
WITH latest_label AS (
	SELECT bundle_id, url
	FROM bundle_labels
	WHERE bundle_id = $1
	ORDER BY created_at DESC
)
SELECT b.*, l.url
FROM bundle b
LEFT JOIN latest_label l ON b.id = l.bundle_id
WHERE b.id = $1
`, id)
	return &bundle, err
}
