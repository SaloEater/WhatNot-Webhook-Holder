CREATE TABLE stream_shipment_space (
    id SERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL,
    shipping_list_id BIGINT NOT NULL
);

ALTER TABLE stream_shipment ADD COLUMN stream_shipment_space_id BIGINT REFERENCES stream_shipment_space(id);

ALTER TABLE stream_shipment DROP COLUMN is_disabled;
ALTER TABLE stream_shipment DROP COLUMN space_name;
ALTER TABLE stream_shipment DROP COLUMN list_id;
ALTER TABLE stream_shipment DROP COLUMN api_key;