CREATE TABLE whatnot_show (
    id             BIGSERIAL PRIMARY KEY,
    channel_id     BIGINT      NOT NULL REFERENCES channel(id),
    show_id        VARCHAR     NOT NULL,
    title          VARCHAR     NOT NULL DEFAULT '',
    description    TEXT        NOT NULL DEFAULT '',
    thumbnail_url  VARCHAR     NOT NULL DEFAULT '',
    visibility     VARCHAR     NOT NULL DEFAULT '',
    scheduled_date TIMESTAMPTZ NULL,
    created_date   TIMESTAMPTZ NOT NULL DEFAULT now(),
    status         VARCHAR     NOT NULL DEFAULT 'planned',
    started_at     TIMESTAMPTZ NULL,
    ended_at       TIMESTAMPTZ NULL,
    CONSTRAINT whatnot_show_channel_show_uniq UNIQUE (channel_id, show_id)
);

CREATE TABLE stream_whatnot_show (
    stream_id       BIGINT      NOT NULL REFERENCES stream(id),
    whatnot_show_id BIGINT      NOT NULL REFERENCES whatnot_show(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (stream_id, whatnot_show_id)
);
CREATE INDEX stream_whatnot_show_show_idx ON stream_whatnot_show (whatnot_show_id);

-- Backfill every stream that has a show id (deleted streams included). When several streams share a show,
-- the newest stream supplies the name.
INSERT INTO whatnot_show (channel_id, show_id, title, status, created_date)
SELECT DISTINCT ON (s.channel_id, s.whatnot_show_id)
    s.channel_id,
    s.whatnot_show_id,
    COALESCE(s.whatnot_show_name, ''),
    CASE s.status WHEN 'live' THEN 'live' WHEN 'ended' THEN 'ended' WHEN 'verified' THEN 'ended' ELSE 'planned' END,
    s.created_at
FROM stream s
WHERE s.whatnot_show_id IS NOT NULL AND s.whatnot_show_id <> ''
ORDER BY s.channel_id, s.whatnot_show_id, s.id DESC;

INSERT INTO stream_whatnot_show (stream_id, whatnot_show_id, created_at)
SELECT s.id, w.id, s.created_at
FROM stream s
INNER JOIN whatnot_show w ON w.channel_id = s.channel_id AND w.show_id = s.whatnot_show_id
WHERE s.whatnot_show_id IS NOT NULL AND s.whatnot_show_id <> '';

ALTER TABLE stream DROP COLUMN whatnot_show_id;
ALTER TABLE stream DROP COLUMN whatnot_show_name;
