CREATE TABLE widget_board_price_ranges (
    id         SERIAL      PRIMARY KEY,
    channel_id BIGINT      NOT NULL,
    tier_id    VARCHAR(32) NOT NULL,
    price_from INTEGER     NOT NULL,
    UNIQUE (channel_id, tier_id)
);
