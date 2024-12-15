-- name: CreateUserScope :one
INSERT INTO user_scopes (user_id, scope_id)
VALUES ($1, $2)
  ON CONFLICT (user_id, scope_id) DO UPDATE SET scope_id = EXCLUDED.scope_id, user_id = EXCLUDED.user_id
RETURNING user_id, scope_id;

-- name: FindScopeByName :one
SELECT id, name FROM scopes WHERE name = $1;
