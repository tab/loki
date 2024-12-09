-- +goose Up
CREATE INDEX tokens_user_id_idx ON tokens (user_id);
CREATE INDEX tokens_expires_at_idx ON tokens (expires_at);

-- +goose Down
DROP INDEX IF EXISTS tokens_user_id_idx;
DROP INDEX IF EXISTS tokens_expires_at_idx;
