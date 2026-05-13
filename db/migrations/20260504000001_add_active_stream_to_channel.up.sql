ALTER TABLE channel ADD COLUMN active_stream_id INT;
ALTER TABLE channel ADD CONSTRAINT fk_active_stream FOREIGN KEY (active_stream_id) REFERENCES stream (id);
