CREATE TABLE channel (
    id SERIAL PRIMARY KEY,
    name varchar(255),
    is_deleted bool
);
INSERT INTO channel (name, is_deleted) values ('Mount Olympus Breaks', false)