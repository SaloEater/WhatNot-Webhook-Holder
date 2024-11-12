CREATE TABLE bundle
(
    id          SERIAL PRIMARY KEY,
    location_id INT,
    status      varchar(255) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted  BOOLEAN      NOT NULL DEFAULT FALSE,
    locked_at   TIMESTAMP,
    name        VARCHAR(255) NOT NULL
);

CREATE TABLE location
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE box_type
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(255) NOT NULL,
    image VARCHAR(255)
);

CREATE TABLE box
(
    id        SERIAL PRIMARY KEY,
    bundle_id INT NOT NULL,
    status    INT NOT NULL,
    boxes_id  INT NOT NULL,
    label_id  INT NOT NULL,
    channel_id INT,
    index     INT NOT NULL
);

CREATE TABLE bundle_boxes
(
    id          SERIAL PRIMARY KEY,
    bundle_id   INT NOT NULL,
    box_type_id INT NOT NULL,
    count       INT NOT NULL
);

CREATE TABLE bundle_labels
(
    id         SERIAL PRIMARY KEY,
    bundle_id  INT       NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    url        VARCHAR(255)
);
ALTER TABLE bundle ADD CONSTRAINT fk_location_id FOREIGN KEY (location_id) REFERENCES location (id);
ALTER TABLE box ADD CONSTRAINT fk_bundle_id FOREIGN KEY (bundle_id) REFERENCES bundle (id);
ALTER TABLE box ADD CONSTRAINT fk_boxes_id FOREIGN KEY (boxes_id) REFERENCES bundle_boxes (id);
ALTER TABLE box ADD CONSTRAINT fk_label_id FOREIGN KEY (label_id) REFERENCES bundle_labels (id);
ALTER TABLE bundle_boxes ADD CONSTRAINT fk_bundle_id FOREIGN KEY (bundle_id) REFERENCES bundle (id);
ALTER TABLE bundle_boxes ADD CONSTRAINT fk_box_type_id FOREIGN KEY (box_type_id) REFERENCES box_type (id);
ALTER TABLE bundle_labels ADD CONSTRAINT fk_bundle_id FOREIGN KEY (bundle_id) REFERENCES bundle (id);