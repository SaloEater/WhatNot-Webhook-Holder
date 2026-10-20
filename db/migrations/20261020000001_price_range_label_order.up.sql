ALTER TABLE series_price_ranges
    ADD COLUMN label      TEXT    NOT NULL DEFAULT '',
    ADD COLUMN sort_order INTEGER NOT NULL DEFAULT 0;
-- Backfill: today's implicit order (price_from, id) becomes the explicit one.
UPDATE series_price_ranges r
SET sort_order = s.rn
FROM (SELECT id, ROW_NUMBER() OVER (PARTITION BY series_id ORDER BY price_from, id) - 1 AS rn
      FROM series_price_ranges) s
WHERE r.id = s.id;
CREATE INDEX series_price_ranges_series_order_idx ON series_price_ranges(series_id, sort_order);
