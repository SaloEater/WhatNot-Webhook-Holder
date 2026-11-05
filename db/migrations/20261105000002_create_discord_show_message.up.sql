CREATE TABLE discord_show_message (
    whatnot_show_id    BIGINT      PRIMARY KEY REFERENCES whatnot_show(id),
    view               VARCHAR     NOT NULL DEFAULT 'overview',
    discord_channel_id VARCHAR     NULL,
    discord_message_id VARCHAR     NULL,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);
