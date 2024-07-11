ALTER TABLE channel ADD COLUMN default_high_bid_team varchar(255) DEFAULT '' NOT NULL;
ALTER TABLE channel ADD COLUMN default_high_bid_floor integer DEFAULT 0 NOT NULL;