-- name: CreateUserRole :one
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
  ON CONFLICT (user_id, role_id) DO UPDATE SET role_id = EXCLUDED.role_id, user_id = EXCLUDED.user_id
RETURNING user_id, role_id;

-- name: FindRoleByName :one
SELECT id, name FROM roles WHERE name = $1;
