CREATE TABLE stream_shipment (
    id SERIAL PRIMARY KEY,
    stream_id BIGINT NOT NULL,
    is_disabled BOOLEAN DEFAULT FALSE,
    api_key VARCHAR(255) NOT NULL DEFAULT '',
    space_name VARCHAR(255) NOT NULL,
    list_id bigint NOT NULL,
    task_id varchar(255) DEFAULT '',
    last_task_status VARCHAR(255) DEFAULT ''
);