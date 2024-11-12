package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
	"strings"
)

type BoxRepository struct {
	DB *sqlx.DB
}

func (r *BoxRepository) GetByBundleIDAndItemIndex(bundleID int64, itemIndex int) (*entity.Box, error) {
	var box entity.Box
	err := r.DB.Get(&box, `SELECT * FROM box WHERE bundle_id = $1 AND item_index = $2`, bundleID, itemIndex)
	return &box, err
}

func (r *BoxRepository) Update(box *entity.Box) error {
	_, err := r.DB.NamedExec(`UPDATE box SET status = :status WHERE id = :id`, box)
	return err
}

func (r *BoxRepository) BatchCreate(boxes []*entity.Box, labelID int64) error {
	const batchSize = 500
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	for i := 0; i < len(boxes); i += batchSize {
		end := i + batchSize
		if end > len(boxes) {
			end = len(boxes)
		}

		query := `INSERT INTO box (bundle_id, status, boxes_id, label_id, index) VALUES `
		values := []string{}
		args := []interface{}{}

		for _, box := range boxes[i:end] {
			values = append(values, "(?, ?, ?, ?, ?)")
			args = append(args, box.BundleID, box.Status, box.BoxesID, labelID, box.Index)
		}

		query += strings.Join(values, ",")
		query = r.DB.Rebind(query)
		stmt, err := tx.Preparex(query)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
