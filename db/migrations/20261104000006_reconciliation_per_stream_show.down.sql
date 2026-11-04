DROP INDEX IF EXISTS show_reconciliation_stream_show_uniq;
DELETE FROM show_reconciliation a USING show_reconciliation b WHERE a.stream_id = b.stream_id AND a.id < b.id;
CREATE UNIQUE INDEX show_reconciliation_stream_uniq ON show_reconciliation (stream_id);
