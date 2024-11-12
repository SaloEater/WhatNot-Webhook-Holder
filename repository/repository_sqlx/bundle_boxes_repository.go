package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type BundleBoxesRepository struct {
	DB *sqlx.DB
}

func (r *BundleBoxesRepository) Create(bundle *entity.BundleBoxes) error {
	var id int64
	err := r.DB.QueryRow(
		`INSERT INTO bundle_boxes (bundle_id, box_type_id, count) VALUES ($1, $2, $3) RETURNING id`,
		bundle.BundleID, bundle.BoxTypeID, bundle.Count,
	).Scan(&id)
	if err == nil {
		bundle.ID = id
	}
	return err
}

func (r *BundleBoxesRepository) GetAllByBundle(bundleID int64) ([]*entity.BundleBoxes, error) {
	var bundles []*entity.BundleBoxes
	err := r.DB.Select(&bundles, `SELECT * FROM bundle_boxes WHERE bundle_id = $1`, bundleID)
	return bundles, err
}

func (r *BundleBoxesRepository) GetByID(id int64) (*entity.BundleBoxes, error) {
	var bundle entity.BundleBoxes
	err := r.DB.Get(&bundle, `SELECT * FROM bundle_boxes WHERE id = $1`, id)
	return &bundle, err
}

func (r *BundleBoxesRepository) Update(bundle *entity.BundleBoxes) error {
	_, err := r.DB.Exec(
		`UPDATE bundle_boxes SET box_type_id = $1, count = $2 WHERE id = $3`,
		bundle.BoxTypeID, bundle.Count, bundle.ID,
	)
	return err
}

func (r *BundleBoxesRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM bundle_boxes WHERE id = $1`, id)
	return err
}

func (r *BundleBoxesRepository) GetByBundleIDAndBoxTypeID(bundleID, boxTypeID int64) (*entity.BundleBoxes, error) {
	var bundle entity.BundleBoxes
	err := r.DB.Get(&bundle, `SELECT * FROM bundle_boxes WHERE bundle_id = $1 AND box_type_id = $2`, bundleID, boxTypeID)
	return &bundle, err
}
