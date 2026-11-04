package repository_sqlx

import (
	"database/sql"
	"errors"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type WhatnotWebhookEventRepository struct {
	DB *sqlx.DB
}

func (r *WhatnotWebhookEventRepository) InsertWebhookEvent(e *entity.WhatnotWebhookEvent) (bool, error) {
	var id int64
	err := r.DB.Get(&id, `INSERT INTO whatnot_webhook_event (webhook_id, channel_id, type, seller_id, event_time, payload)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (webhook_id) DO NOTHING
RETURNING id`, e.WebhookId, e.ChannelId, e.Type, e.SellerId, e.EventTime, []byte(e.Payload))
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	e.Id = id
	return true, nil
}

func (r *WhatnotWebhookEventRepository) SetWebhookEventStatus(id int64, status string, errMsg string) error {
	_, err := r.DB.Exec(`UPDATE whatnot_webhook_event SET status = $2, error = $3, processed_at = now() WHERE id = $1`, id, status, errMsg)
	return err
}

func (r *WhatnotWebhookEventRepository) GetWebhookEvent(id int64) (*entity.WhatnotWebhookEvent, error) {
	var e entity.WhatnotWebhookEvent
	err := r.DB.Unsafe().Get(&e, `SELECT * FROM whatnot_webhook_event WHERE id = $1`, id)
	return &e, err
}

func (r *WhatnotWebhookEventRepository) ListWebhookEventsByChannel(channelId int64, limit int) ([]*entity.WhatnotWebhookEvent, error) {
	events := []*entity.WhatnotWebhookEvent{}
	err := r.DB.Unsafe().Select(&events, `SELECT id, webhook_id, channel_id, type, seller_id, event_time, received_at, status, error, processed_at
FROM whatnot_webhook_event WHERE channel_id = $1 ORDER BY received_at DESC, id DESC LIMIT $2`, channelId, limit)
	return events, err
}
