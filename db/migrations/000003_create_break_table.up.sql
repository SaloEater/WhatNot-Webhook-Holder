CREATE TABLE break (
    id SERIAL PRIMARY KEY,
    day_id INT,
    name varchar(255),
    start_date DATE,
    end_date DATE,
    CONSTRAINT fk_day
        FOREIGN KEY (day_id)
           REFERENCES day(id)
)