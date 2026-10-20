DROP INDEX IF EXISTS series_price_ranges_series_order_idx;
ALTER TABLE series_price_ranges DROP COLUMN IF EXISTS sort_order, DROP COLUMN IF EXISTS label;
