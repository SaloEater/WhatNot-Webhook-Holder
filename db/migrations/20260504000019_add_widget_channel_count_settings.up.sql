CREATE TABLE widget_channel_count_settings (
    id              SERIAL  PRIMARY KEY,
    channel_id      BIGINT  NOT NULL UNIQUE,
    show_percentage BOOLEAN NOT NULL DEFAULT false
);
