-- name: CreateUserRole :one
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
  ON CONFLICT (user_id, role_id) DO UPDATE SET role_id = EXCLUDED.role_id, user_id = EXCLUDED.user_id
RETURNING user_id, role_id;

-- name: FindRoleByName :one
SELECT id, name FROM roles WHERE name = $1;

-- name: FindUserRoles :many
SELECT id, name FROM roles WHERE id IN (SELECT role_id FROM user_roles WHERE user_id = $1);
