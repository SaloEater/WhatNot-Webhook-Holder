CREATE TABLE widget_presets (
    id         SERIAL      PRIMARY KEY,
    channel_id BIGINT      NOT NULL,
    name       VARCHAR(64) NOT NULL,
    value      TEXT        NOT NULL,
    UNIQUE (channel_id, name)
);
