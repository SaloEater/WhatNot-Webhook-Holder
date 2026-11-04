CREATE TABLE show_reconciliation (
    id           BIGSERIAL PRIMARY KEY,
    stream_id    BIGINT      NOT NULL REFERENCES stream(id),
    show_id      VARCHAR(64) NOT NULL,
    trigger      VARCHAR(16) NOT NULL,
    status       VARCHAR(16) NOT NULL DEFAULT 'running',
    started_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at  TIMESTAMPTZ NULL,
    error        TEXT        NOT NULL DEFAULT '',
    orders       JSONB       NOT NULL DEFAULT '{}'::jsonb,
    report       JSONB       NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX show_reconciliation_stream_idx ON show_reconciliation (stream_id, started_at DESC);
