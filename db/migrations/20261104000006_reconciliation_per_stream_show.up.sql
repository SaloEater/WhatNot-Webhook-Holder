-- One reconciliation per stream + show: linking a different show starts a new row from scratch.
DROP INDEX IF EXISTS show_reconciliation_stream_uniq;
CREATE UNIQUE INDEX show_reconciliation_stream_show_uniq ON show_reconciliation (stream_id, show_id);
