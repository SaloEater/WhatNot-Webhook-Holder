CREATE TABLE series_team_price (
    id        BIGSERIAL PRIMARY KEY,
    series_id BIGINT    NOT NULL REFERENCES series(id),
    team      VARCHAR   NOT NULL,
    price     NUMERIC   NOT NULL DEFAULT 0,
    UNIQUE (series_id, team)
);
