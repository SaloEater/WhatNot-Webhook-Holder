CREATE TABLE channel (
    id SERIAL PRIMARY KEY,
    name varchar(255),
    is_deleted bool
);
INSERT INTO channel (name) values ('Mount Olympus Breaks')