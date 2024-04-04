package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type BreakRepository struct {
	DB *sqlx.DB
}

func (r *BreakRepository) Create(*entity.Break) (int, error) {
	return 0, nil
}

func (r *BreakRepository) Get(int) (*entity.Break, error) {
	return nil, nil
}

func (r *BreakRepository) Delete(int) error {
	return nil
}

func (r *BreakRepository) Update(p *entity.Break) error {
	return nil
}
