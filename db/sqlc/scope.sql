-- name: CreateUserScope :one
INSERT INTO user_scopes (user_id, scope_id)
VALUES ($1, $2)
  ON CONFLICT (user_id, scope_id) DO UPDATE SET scope_id = EXCLUDED.scope_id, user_id = EXCLUDED.user_id
RETURNING user_id, scope_id;

-- name: FindScopeByName :one
SELECT id, name FROM scopes WHERE name = $1;

-- name: UpsertUserScopeByName :one
INSERT INTO user_scopes (user_id, scope_id)
VALUES ($1, (SELECT id FROM scopes WHERE name = $2))
  ON CONFLICT (user_id, scope_id) DO UPDATE SET scope_id = EXCLUDED.scope_id, user_id = EXCLUDED.user_id
RETURNING user_id, scope_id;

-- name: FindUserScopes :many
SELECT id, name FROM scopes WHERE id IN (SELECT scope_id FROM user_scopes WHERE user_id = $1);
