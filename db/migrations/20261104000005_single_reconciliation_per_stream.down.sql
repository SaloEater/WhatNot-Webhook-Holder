ALTER TABLE show_reconciliation DROP COLUMN manual_map;
DROP INDEX IF EXISTS show_reconciliation_stream_uniq;
CREATE INDEX show_reconciliation_stream_idx ON show_reconciliation (stream_id, started_at DESC);
