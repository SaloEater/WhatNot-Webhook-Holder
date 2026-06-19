CREATE TABLE widget_series_boxes_per_break (
    id        SERIAL  PRIMARY KEY,
    series_id BIGINT  NOT NULL UNIQUE,
    amount    INTEGER NOT NULL DEFAULT 0
);
