ALTER TABLE stream ADD COLUMN whatnot_show_id VARCHAR(64) NULL;
ALTER TABLE stream ADD COLUMN whatnot_show_name TEXT NULL;

-- A stream gets back its most recently linked show.
UPDATE stream s
SET whatnot_show_id = x.show_id,
    whatnot_show_name = NULLIF(x.title, '')
FROM (
    SELECT DISTINCT ON (l.stream_id) l.stream_id, w.show_id, w.title
    FROM stream_whatnot_show l
    INNER JOIN whatnot_show w ON w.id = l.whatnot_show_id
    ORDER BY l.stream_id, l.created_at DESC, w.id DESC
) x
WHERE s.id = x.stream_id;

DROP TABLE stream_whatnot_show;
DROP TABLE whatnot_show;
