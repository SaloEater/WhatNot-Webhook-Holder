CREATE TABLE photo (
    id         BIGSERIAL PRIMARY KEY,
    series_id  BIGINT      NOT NULL REFERENCES series(id),
    name       TEXT        NOT NULL DEFAULT '',
    url        TEXT        NOT NULL,
    is_sold    BOOLEAN     NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_deleted BOOLEAN     NOT NULL DEFAULT false
);
