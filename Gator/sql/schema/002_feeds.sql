-- +goose Up
CREATE TABLE feeds (
  id UUID UNIQUE,
  name TEXT,
  url TEXT UNIQUE,
  updated_at TIMESTAMP,
  user_id UUID,
  CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
