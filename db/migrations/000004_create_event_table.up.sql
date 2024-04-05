CREATE TABLE event (
    id SERIAL PRIMARY KEY,
    break_id INT,
    index INT,
    customer varchar(255),
    price float8,
    team varchar(255),
    is_giveaway bit,
    note text,
    quantity int,
    CONSTRAINT fk_break
       FOREIGN KEY (break_id)
           REFERENCES break(id)
)