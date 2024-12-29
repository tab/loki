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
