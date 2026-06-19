CREATE TABLE widget_series_stashorpass (
    id         SERIAL PRIMARY KEY,
    channel_id bigint NOT NULL UNIQUE,
    price      integer NOT NULL DEFAULT 0
);

CREATE TABLE widget_series_pick2 (
    id         SERIAL PRIMARY KEY,
    channel_id bigint NOT NULL UNIQUE,
    price      integer NOT NULL DEFAULT 0
);
