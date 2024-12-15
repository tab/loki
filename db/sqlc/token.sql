-- name: CreateToken :one
INSERT INTO tokens (user_id, type, value, expires_at)
VALUES ($1, $2, $3, $4)
  RETURNING id, type, value, expires_at;
