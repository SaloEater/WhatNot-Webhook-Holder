package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type WidgetBoardPriceRangeRepository struct {
	DB *sqlx.DB
}

func (r *WidgetBoardPriceRangeRepository) ListByChannel(channelId int64) ([]*entity.WidgetBoardPriceRange, error) {
	var rows []*entity.WidgetBoardPriceRange
	err := r.DB.Unsafe().Select(&rows, `SELECT * FROM widget_board_price_ranges WHERE channel_id = $1 ORDER BY price_from DESC`, channelId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *WidgetBoardPriceRangeRepository) Create(w *entity.WidgetBoardPriceRange) (int64, error) {
	var id int64
	err := r.DB.QueryRow(
		`INSERT INTO widget_board_price_ranges (channel_id, tier_id, price_from) VALUES ($1, $2, $3) RETURNING id`,
		w.ChannelId, w.TierId, w.PriceFrom,
	).Scan(&id)
	return id, err
}

func (r *WidgetBoardPriceRangeRepository) Upsert(channelId int64, tierId string, priceFrom int) error {
	_, err := r.DB.Exec(`
		INSERT INTO widget_board_price_ranges (channel_id, tier_id, price_from)
		VALUES ($1, $2, $3)
		ON CONFLICT (channel_id, tier_id) DO UPDATE SET price_from = $3`,
		channelId, tierId, priceFrom)
	return err
}

func (r *WidgetBoardPriceRangeRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM widget_board_price_ranges WHERE id=$1`, id)
	return err
}
