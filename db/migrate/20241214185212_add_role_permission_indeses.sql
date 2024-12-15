-- +goose Up
CREATE INDEX role_permissions_role_id_idx ON role_permissions (role_id);
CREATE INDEX role_permissions_permission_id_idx ON role_permissions (permission_id);

-- +goose Down
DROP INDEX IF EXISTS role_permissions_role_id_idx;
DROP INDEX IF EXISTS role_permissions_permission_id_idx;
