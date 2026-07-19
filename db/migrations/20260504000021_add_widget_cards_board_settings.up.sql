CREATE TABLE widget_cards_board_settings (
    id          SERIAL      PRIMARY KEY,
    channel_id  BIGINT      NOT NULL UNIQUE,
    orientation VARCHAR(16) NOT NULL DEFAULT 'list'
);
