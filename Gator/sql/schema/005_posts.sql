-- +goose Up
CREATE TABLE posts (
  id UUID UNIQUE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  published_at TIMESTAMP,
  title TEXT,
  url TEXT UNIQUE,
  description TEXT,
  feed_id UUID,
  CONSTRAINT fk_feed_id
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
