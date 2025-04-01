-- name: FindScopes :many
WITH counter AS (
  SELECT COUNT(*) AS total
  FROM scopes
)
SELECT
  s.id,
  s.name,
  s.description,
  counter.total
FROM scopes AS s
RIGHT JOIN counter ON TRUE
ORDER BY s.created_at DESC LIMIT $1::bigint OFFSET $2::bigint;

-- name: CreateScope :one
INSERT INTO scopes (name, description)
VALUES ($1, $2)
  RETURNING id, name, description;

-- name: FindScopeById :one
SELECT id, name, description FROM scopes WHERE id = $1;

-- name: FindScopeByName :one
SELECT id, name FROM scopes WHERE name = $1;

-- name: UpdateScope :one
UPDATE scopes
SET
  name = $2,
  description = $3,
  updated_at = NOW()
WHERE id = $1
RETURNING id, name, description;

-- name: DeleteScope :exec
DELETE FROM scopes WHERE id = $1;

-- name: CreateUserScope :one
INSERT INTO user_scopes (user_id, scope_id)
VALUES ($1, $2)
  ON CONFLICT (user_id, scope_id) DO UPDATE SET scope_id = EXCLUDED.scope_id, user_id = EXCLUDED.user_id
RETURNING user_id, scope_id;

-- name: CreateUserScopes :many
WITH
  deleted AS (
    DELETE FROM user_scopes
    WHERE user_id = @user_id::uuid AND scope_id NOT IN (SELECT unnest(@scope_ids::uuid[]))
  ),
  inserted AS (
    INSERT INTO user_scopes (user_id, scope_id)
    SELECT @user_id::uuid, scope_id
    FROM unnest(@scope_ids::uuid[]) AS scope_id
    ON CONFLICT (user_id, scope_id) DO NOTHING
      RETURNING user_id, scope_id
  )
SELECT user_id, scope_id FROM inserted;

-- name: UpsertUserScopeByName :one
INSERT INTO user_scopes (user_id, scope_id)
VALUES ($1, (SELECT id FROM scopes WHERE name = $2))
  ON CONFLICT (user_id, scope_id) DO UPDATE SET scope_id = EXCLUDED.scope_id, user_id = EXCLUDED.user_id
RETURNING user_id, scope_id;

-- name: FindUserScopes :many
SELECT id, name FROM scopes WHERE id IN (SELECT scope_id FROM user_scopes WHERE user_id = $1);
