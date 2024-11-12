package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type BoxTypeRepository struct {
	DB *sqlx.DB
}

func (r *BoxTypeRepository) GetAll() ([]*entity.BoxType, error) {
	bundleBoxes := make([]*entity.BoxType, 0)
	err := r.DB.Select(&bundleBoxes, "SELECT * FROM box_type")
	if err != nil {
		return nil, err
	}
	return bundleBoxes, nil
}

func (r *BoxTypeRepository) GetByID(id int64) (*entity.BoxType, error) {
	boxType := &entity.BoxType{}
	err := r.DB.Get(boxType, "SELECT * FROM box_type WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return boxType, nil
}

func (r *BoxTypeRepository) Create(boxType *entity.BoxType) error {
	var id int64
	err := r.DB.QueryRow("INSERT INTO box_type (name, image) VALUES ($1, $2) RETURNING id", boxType.Name, boxType.Image).Scan(&id)
	if err != nil {
		return err
	}
	boxType.ID = id
	return nil
}

func (r *BoxTypeRepository) Update(boxType *entity.BoxType) error {
	_, err := r.DB.Exec("UPDATE box_type SET name = $1, image = $2 WHERE id = $3", boxType.Name, boxType.Image, boxType.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *BoxTypeRepository) GetByIDs(ids []int64) ([]*entity.BoxType, error) {
	boxTypes := make([]*entity.BoxType, 0)
	query, args, err := sqlx.In("SELECT * FROM box_type WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}
	err = r.DB.Select(&boxTypes, r.DB.Rebind(query), args...)
	if err != nil {
		return nil, err
	}
	return boxTypes, nil
}
