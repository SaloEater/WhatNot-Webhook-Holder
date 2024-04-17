CREATE TABLE demo (
    id SERIAL PRIMARY KEY,
    stream_id INT NOT NULL,
    highlight_username varchar(255) NOT NULL DEFAULT '',
    break_id INT NOT NULL DEFAULT 0,
    CONSTRAINT fk_stream
        FOREIGN KEY (stream_id)
            REFERENCES stream(id),
    CONSTRAINT fk_break
        FOREIGN KEY (break_id)
              REFERENCES break(id),
    CONSTRAINT unique_demo
        UNIQUE (stream_id)
)