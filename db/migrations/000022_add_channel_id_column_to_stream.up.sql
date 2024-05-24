ALTER TABLE stream ADD COLUMN channel_id INT NOT NULL DEFAULT 0;
ALTER TABLE stream ADD CONSTRAINT fk_channel_id FOREIGN KEY (channel_id) REFERENCES channel(id);