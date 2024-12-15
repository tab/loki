-- +goose Up
CREATE INDEX user_scopes_user_id_idx ON user_scopes (user_id);
CREATE INDEX user_scopes_scope_id_idx ON user_scopes (scope_id);

-- +goose Down
DROP INDEX IF EXISTS user_scopes_user_id_idx;
DROP INDEX IF EXISTS user_scopes_scope_id_idx;
