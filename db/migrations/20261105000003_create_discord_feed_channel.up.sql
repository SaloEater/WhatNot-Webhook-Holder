CREATE TABLE discord_feed_channel (
    channel_id         BIGINT      PRIMARY KEY REFERENCES channel(id),
    discord_channel_id VARCHAR     NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);
