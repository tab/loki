-- name: FindUserPermissions :many
SELECT id, name FROM permissions WHERE id IN (
  SELECT permission_id FROM role_permissions WHERE role_id IN (
    SELECT role_id FROM user_roles WHERE user_id = $1));
