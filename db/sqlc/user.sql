-- name: FindUsers :many
WITH counter AS (
  SELECT COUNT(*) AS total
  FROM users
)
SELECT
  u.id,
  u.identity_number,
  u.personal_code,
  u.first_name,
  u.last_name,
  counter.total
FROM users AS u
RIGHT JOIN counter ON TRUE
ORDER BY u.created_at DESC LIMIT $1::bigint OFFSET $2::bigint;

-- name: CreateUser :one
INSERT INTO users (identity_number, personal_code, first_name, last_name)
VALUES ($1, $2, $3, $4)
  ON CONFLICT (identity_number) DO UPDATE SET first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name
RETURNING id, identity_number, personal_code, first_name, last_name;

-- name: FindUserById :one
SELECT id, identity_number, personal_code, first_name, last_name FROM users WHERE id = $1;

-- name: FindUserByIdentityNumber :one
SELECT id, identity_number, personal_code, first_name, last_name FROM users WHERE identity_number = $1;

-- name: FindUserDetailsById :one
SELECT
  u.id,
  u.identity_number,
  u.personal_code,
  u.first_name,
  u.last_name,
  COALESCE(ur.roles, ARRAY[]::uuid[]) AS role_ids,
  COALESCE(us.scopes, ARRAY[]::uuid[]) AS scope_ids
FROM users u
  LEFT JOIN (
    SELECT
      ur.user_id,
      ARRAY_AGG(ur.role_id) AS roles
    FROM user_roles ur
    GROUP BY ur.user_id
  ) ur ON u.id = ur.user_id
  LEFT JOIN (
    SELECT
      us.user_id,
      ARRAY_AGG(us.scope_id) AS scopes
    FROM user_scopes us
    GROUP BY us.user_id
  ) us ON u.id = us.user_id
WHERE
  u.id = $1;

-- name: UpdateUser :one
UPDATE users
SET
  identity_number = $2,
  personal_code = $3,
  first_name = $4,
  last_name = $5,
  updated_at = NOW()
WHERE id = $1
RETURNING id, identity_number, personal_code, first_name, last_name;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
