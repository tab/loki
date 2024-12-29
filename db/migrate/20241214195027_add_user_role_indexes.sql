-- +goose Up
CREATE INDEX user_roles_user_id_idx ON user_roles (user_id);
CREATE INDEX user_roles_role_id_idx ON user_roles (role_id);

-- +goose Down
DROP INDEX user_roles_user_id_idx;
DROP INDEX user_roles_role_id_idx;

