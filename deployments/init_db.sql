CREATE TABLE images
(
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    tg_file_id VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    usage_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, tg_file_id)
);

CREATE TABLE tags
(
    image_id INT NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    position INT NOT NULL,
    PRIMARY KEY(image_id, name)
);

CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_tags_image_id ON tags(image_id);
