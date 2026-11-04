ALTER TABLE stream ADD COLUMN status VARCHAR(16) NOT NULL DEFAULT 'planned';
UPDATE stream SET status = 'ended' WHERE is_ended = true;
