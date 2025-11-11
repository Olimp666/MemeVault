CREATE TABLE images
(
    tg_file_id VARCHAR(255) PRIMARY KEY,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tags
(
    tg_file_id VARCHAR(255) NOT NULL REFERENCES images(tg_file_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY(tg_file_id, name)
);

CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_tags_tg_file_id ON tags(tg_file_id);
