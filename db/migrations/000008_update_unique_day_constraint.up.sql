ALTER TABLE day DROP CONSTRAINT unique_date;
CREATE UNIQUE INDEX idx_unique_date_excluding_deleted ON day (date) WHERE NOT is_deleted;