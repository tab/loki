-- name: CreateUser :one
INSERT INTO users (identity_number, personal_code, first_name, last_name)
VALUES ($1, $2, $3, $4)
  ON CONFLICT (identity_number) DO UPDATE SET first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name
RETURNING id, identity_number, personal_code, first_name, last_name;

-- name: FindUserByIdentityNumber :one
SELECT id, identity_number, personal_code, first_name, last_name FROM users WHERE identity_number = $1;

-- name: FindUserById :one
SELECT id, identity_number, personal_code, first_name, last_name FROM users WHERE id = $1;

