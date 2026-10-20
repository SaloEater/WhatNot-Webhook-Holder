CREATE TABLE series_preset (
    id         BIGSERIAL   PRIMARY KEY,
    channel_id BIGINT      NOT NULL,
    name       TEXT        NOT NULL,
    series_ids JSONB       NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (channel_id, name)
);
