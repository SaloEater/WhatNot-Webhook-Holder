ALTER TABLE stream ADD COLUMN active_break_id INT;
ALTER TABLE stream ADD COLUMN highlight_username VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE stream ADD CONSTRAINT fk_active_break FOREIGN KEY (active_break_id) REFERENCES break (id);
