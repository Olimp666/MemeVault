CREATE TABLE images
(
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    data BYTEA NOT NULL
);

CREATE TABLE tags
(
    id SERIAL PRIMARY KEY,
    image_id BIGINT NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    UNIQUE(image_id, name)
);

CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_tags_image_id ON tags(image_id);
