CREATE TABLE tg_chat (
                         id BIGINT PRIMARY KEY,
                         chat_id BIGINT NOT NULL,
                         is_disabled BOOLEAN NOT NULL DEFAULT FALSE,
                         added_date TIMESTAMP NOT NULL DEFAULT NOW()
)