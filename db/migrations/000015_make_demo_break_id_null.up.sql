ALTER TABLE demo DROP COLUMN break_id;
ALTER TABLE demo ADD COLUMN break_id int default null;
ALTER TABLE demo ADD CONSTRAINT fk_break foreign key(break_id) references break(id);