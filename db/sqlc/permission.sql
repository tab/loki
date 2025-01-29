-- name: FindPermissions :many
WITH counter AS (
  SELECT COUNT(*) AS total
  FROM permissions
)
SELECT
  p.id,
  p.name,
  p.description,
  counter.total
FROM permissions AS p
RIGHT JOIN counter ON TRUE
ORDER BY p.created_at DESC LIMIT $1 OFFSET $2;

-- name: CreatePermission :one
INSERT INTO permissions (name, description)
VALUES ($1, $2)
RETURNING id, name, description;

-- name: FindPermissionById :one
SELECT id, name, description FROM permissions WHERE id = $1;

-- name: UpdatePermission :one
UPDATE permissions
SET
  name = $2,
  description = $3,
  updated_at = NOW()
WHERE id = $1
RETURNING id, name, description;

-- name: DeletePermission :exec
DELETE FROM permissions WHERE id = $1;

-- name: CreateRolePermissions :many
WITH
  deleted AS (
    DELETE FROM role_permissions
    WHERE role_id = @role_id::uuid AND permission_id NOT IN (SELECT unnest(@permission_ids::uuid[]))
  ),
  inserted AS (
    INSERT INTO role_permissions (role_id, permission_id)
    SELECT @role_id::uuid, permission_id
    FROM unnest(@permission_ids::uuid[]) AS permission_id
    ON CONFLICT (role_id, permission_id) DO NOTHING
      RETURNING role_id, permission_id
  )
SELECT role_id, permission_id FROM inserted;

-- name: FindUserPermissions :many
SELECT id, name FROM permissions WHERE id IN (
  SELECT permission_id FROM role_permissions WHERE role_id IN (
    SELECT role_id FROM user_roles WHERE user_id = $1));
