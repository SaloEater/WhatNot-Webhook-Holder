package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type TrackingRepository struct {
	DB *sqlx.DB
}

func (r *TrackingRepository) GetValidationData(data *entity.ValidationDataBuilder) (*entity.ValidationData, error) {
	var aggData entity.ValidationData

	query := `
		SELECT 
			l.id AS location_id, 
			b.id AS bundle_id, 
			bt.id AS box_type_id, 
			bx.index AS box_index, 
			bx.id AS box_id, 
			bx.status AS box_status, 
			b.location_id AS bundle_loc_id
		FROM 
			location l
		JOIN 
			bundle b ON b.id = $1
		JOIN 
			box_type bt ON bt.id = $2
		JOIN 
			bundle_boxes bb ON bb.bundle_id = b.id AND bb.box_type_id = bt.id
		JOIN 
			box bx ON bx.bundle_id = b.id AND bx.index = $3
		WHERE 
			l.id = $4
	`

	err := r.DB.QueryRow(query, data.BundleID, data.BoxTypeID, data.BoxIndex, data.LocationID).Scan(
		&aggData.LocationID,
		&aggData.BundleID,
		&aggData.BoxTypeID,
		&aggData.BoxIndex,
		&aggData.BoxID,
		&aggData.BoxStatus,
		&aggData.BundleLocID,
	)
	return &aggData, err
}
