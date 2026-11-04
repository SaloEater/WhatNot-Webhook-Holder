-- One reconciliation per stream: re-runs reuse the row. Keeps the newest row of each stream.
DELETE FROM show_reconciliation a USING show_reconciliation b WHERE a.stream_id = b.stream_id AND a.id < b.id;
DROP INDEX IF EXISTS show_reconciliation_stream_idx;
CREATE UNIQUE INDEX show_reconciliation_stream_uniq ON show_reconciliation (stream_id);
-- Human decisions from the manual page, honoured by re-runs: {"breaks":{"<wn name>":<break id or -1 skip>},"giveaways":{"<order id>":<break id>}}
ALTER TABLE show_reconciliation ADD COLUMN manual_map JSONB NOT NULL DEFAULT '{}'::jsonb;
