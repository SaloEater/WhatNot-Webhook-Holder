UPDATE stream s
SET
    active_break_id    = d.break_id,
    highlight_username = d.highlight_username
FROM demo d
WHERE d.stream_id = s.id;

UPDATE channel c
SET active_stream_id = d.stream_id
FROM demo d
WHERE d.id = c.demo_id;
