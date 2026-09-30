-- +goose Up
CREATE TABLE feeds (
  name TEXT,
  url TEXT UNIQUE,
  user_id UUID,
  CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
