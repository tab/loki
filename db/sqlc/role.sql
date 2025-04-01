-- name: FindRoles :many
WITH counter AS (
  SELECT COUNT(*) AS total
  FROM roles
)
SELECT
  r.id,
  r.name,
  r.description,
  counter.total
FROM roles AS r
RIGHT JOIN counter ON TRUE
ORDER BY r.created_at DESC LIMIT $1::bigint OFFSET $2::bigint;

-- name: CreateRole :one
INSERT INTO roles (name, description)
VALUES ($1, $2)
  RETURNING id, name, description;

-- name: FindRoleById :one
SELECT id, name, description FROM roles WHERE id = $1;

-- name: FindRoleByName :one
SELECT id, name FROM roles WHERE name = $1;

-- name: FindRoleDetailsById :one
SELECT
  r.id,
  r.name,
  r.description,
  COALESCE(rp.permissions, ARRAY[]::uuid[]) AS permission_ids
FROM roles r
  LEFT JOIN (
    SELECT
      rp.role_id,
      ARRAY_AGG(rp.permission_id) AS permissions
    FROM role_permissions rp
    GROUP BY rp.role_id
  ) rp ON r.id = rp.role_id
WHERE
  r.id = $1;

-- name: UpdateRole :one
UPDATE roles
SET
  name = $2,
  description = $3,
  updated_at = NOW()
WHERE id = $1
RETURNING id, name, description;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;

-- name: FindUserRoles :many
SELECT id, name FROM roles WHERE id IN (SELECT role_id FROM user_roles WHERE user_id = $1);

-- name: CreateUserRole :one
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
  ON CONFLICT (user_id, role_id) DO UPDATE SET role_id = EXCLUDED.role_id, user_id = EXCLUDED.user_id
RETURNING user_id, role_id;

-- name: CreateUserRoles :many
WITH
  deleted AS (
    DELETE FROM user_roles
    WHERE user_id = @user_id::uuid AND role_id NOT IN (SELECT unnest(@role_ids::uuid[]))
  ),
  inserted AS (
    INSERT INTO user_roles (user_id, role_id)
    SELECT @user_id::uuid, role_id
    FROM unnest(@role_ids::uuid[]) AS role_id
    ON CONFLICT (user_id, role_id) DO NOTHING
      RETURNING user_id, role_id
  )
SELECT user_id, role_id FROM inserted;

-- name: UpsertUserRoleByName :one
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, (SELECT id FROM roles WHERE name = $2))
  ON CONFLICT (user_id, role_id) DO UPDATE SET role_id = EXCLUDED.role_id, user_id = EXCLUDED.user_id
RETURNING user_id, role_id;

