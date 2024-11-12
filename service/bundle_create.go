package service

import (
	"database/sql"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"time"
)

type BundleCreateRequest struct {
	Name string `json:"name"`
}

type BundleResponse struct {
	ID         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	LocationID string    `db:"location_id" json:"location_id"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted  bool      `db:"is_deleted" json:"is_deleted"`
	LockedAt   string    `db:"locked_at" json:"locked_at"`
}

func (s *Service) BundleCreate(request *BundleCreateRequest) (*BundleResponse, error) {
	bundle := &entity.Bundle{
		Status:     entity.BundleStatusPlanned,
		LocationID: sql.NullInt64{Valid: false},
		IsDeleted:  false,
		Name:       request.Name,
	}
	err := s.BundleRepositorier.Create(bundle)
	if err != nil {
		return nil, err
	}
	return mapBundleResponse(bundle), nil
}

func mapBundleResponse(bundle *entity.Bundle) *BundleResponse {
	var locationID string
	if bundle.LocationID.Valid {
		locationID = fmt.Sprintf("%d", bundle.LocationID.Int64)
	} else {
		locationID = ""
	}
	var lockedAt string
	if bundle.LockedAt.Valid {
		lockedAt = bundle.LockedAt.Time.Format(time.RFC3339)
	} else {
		lockedAt = ""
	}
	return &BundleResponse{
		ID:         bundle.ID,
		Name:       bundle.Name,
		LocationID: locationID,
		Status:     bundle.Status,
		CreatedAt:  bundle.CreatedAt,
		UpdatedAt:  bundle.UpdatedAt,
		IsDeleted:  bundle.IsDeleted,
		LockedAt:   lockedAt,
	}
}
