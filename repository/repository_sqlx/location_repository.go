package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type LocationRepository struct {
	DB *sqlx.DB
}

func (r *LocationRepository) GetAll() ([]*entity.Location, error) {
	var locations []*entity.Location
	err := r.DB.Select(&locations, "SELECT * FROM location")
	return locations, err
}

func (r *LocationRepository) GetByID(id int64) (*entity.Location, error) {
	var location entity.Location
	err := r.DB.Get(&location, "SELECT * FROM location WHERE id=$1", id)
	return &location, err
}
