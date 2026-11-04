package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type WhatnotWebhookEventRepositorier interface {
	// InsertWebhookEvent returns inserted=false when the webhook_id already exists (duplicate delivery).
	InsertWebhookEvent(*entity.WhatnotWebhookEvent) (inserted bool, err error)
	SetWebhookEventStatus(id int64, status string, errMsg string) error
	GetWebhookEvent(id int64) (*entity.WhatnotWebhookEvent, error)
	// ListWebhookEventsByChannel returns newest first, without payload.
	ListWebhookEventsByChannel(channelId int64, limit int) ([]*entity.WhatnotWebhookEvent, error)
}
