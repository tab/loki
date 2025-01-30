-- +goose Up
CREATE TABLE user_scopes (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  scope_id UUID NOT NULL REFERENCES scopes(id) ON DELETE CASCADE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, scope_id)
);

-- +goose Down
DROP TABLE user_scopes;
