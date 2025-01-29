-- name: FindTokens :many
WITH counter AS (
  SELECT COUNT(*) AS total
  FROM tokens
)
SELECT
  t.id,
  t.user_id,
  t.type,
  t.value,
  t.expires_at,
  u.identity_number,
  counter.total
FROM tokens AS t
  JOIN users AS u ON t.user_id = u.id
  RIGHT JOIN counter ON TRUE
ORDER BY t.created_at DESC LIMIT $1 OFFSET $2;

-- name: CreateToken :one
INSERT INTO tokens (user_id, type, value, expires_at)
VALUES ($1, $2, $3, $4)
  RETURNING id, type, value, expires_at;

-- name: CreateTokens :many
INSERT INTO tokens (user_id, type, value, expires_at)
VALUES
  (@user_id::uuid, 'access_token'::token_type, @access_token_value::text, @access_token_expires_at::timestamp),
  (@user_id::uuid, 'refresh_token'::token_type, @refresh_token_value::text, @refresh_token_expires_at::timestamp)
  RETURNING id, type, value, expires_at;

-- name: FindTokenById :one
SELECT id, user_id, type, value, expires_at FROM tokens WHERE id = $1;

-- name: DeleteToken :exec
DELETE FROM tokens WHERE id = $1;
