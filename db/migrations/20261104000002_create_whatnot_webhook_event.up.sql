CREATE TABLE whatnot_webhook_event (
    id            BIGSERIAL PRIMARY KEY,
    webhook_id    VARCHAR(64)  NOT NULL UNIQUE,
    channel_id    BIGINT       NOT NULL REFERENCES channel(id),
    type          VARCHAR(64)  NOT NULL,
    seller_id     VARCHAR(32)  NOT NULL DEFAULT '',
    event_time    TIMESTAMPTZ  NULL,
    received_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    payload       JSONB        NOT NULL,
    status        VARCHAR(16)  NOT NULL DEFAULT 'received',
    error         TEXT         NOT NULL DEFAULT '',
    processed_at  TIMESTAMPTZ  NULL
);
CREATE INDEX whatnot_webhook_event_channel_idx ON whatnot_webhook_event (channel_id, received_at DESC);
