package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type BundleLabelRepository struct {
	DB *sqlx.DB
}

func (b *BundleLabelRepository) Create(labels *entity.BundleLabels) error {
	var id int64
	err := b.DB.QueryRow(`INSERT INTO bundle_labels (bundle_id, url, created_at) VALUES ($1, $2, NOW()) RETURNING id`, labels.BundleID, labels.URL).Scan(&id)
	if err != nil {
		return err
	}
	labels.ID = id
	return nil
}
