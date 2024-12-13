// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createToken = `-- name: CreateToken :one
INSERT INTO tokens (user_id, type, value, expires_at)
VALUES ($1, $2, $3, $4)
  RETURNING id, type, value, expires_at
`

type CreateTokenParams struct {
	UserID    uuid.UUID
	Type      TokenType
	Value     string
	ExpiresAt pgtype.Timestamp
}

type CreateTokenRow struct {
	ID        uuid.UUID
	Type      TokenType
	Value     string
	ExpiresAt pgtype.Timestamp
}

func (q *Queries) CreateToken(ctx context.Context, arg CreateTokenParams) (CreateTokenRow, error) {
	row := q.db.QueryRow(ctx, createToken,
		arg.UserID,
		arg.Type,
		arg.Value,
		arg.ExpiresAt,
	)
	var i CreateTokenRow
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Value,
		&i.ExpiresAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (identity_number, personal_code, first_name, last_name)
VALUES ($1, $2, $3, $4)
  ON CONFLICT (identity_number) DO UPDATE SET first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name
RETURNING id, identity_number, personal_code, first_name, last_name
`

type CreateUserParams struct {
	IdentityNumber string
	PersonalCode   string
	FirstName      string
	LastName       string
}

type CreateUserRow struct {
	ID             uuid.UUID
	IdentityNumber string
	PersonalCode   string
	FirstName      string
	LastName       string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.IdentityNumber,
		arg.PersonalCode,
		arg.FirstName,
		arg.LastName,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.IdentityNumber,
		&i.PersonalCode,
		&i.FirstName,
		&i.LastName,
	)
	return i, err
}

const findUserByIdentityNumber = `-- name: FindUserByIdentityNumber :one
SELECT id, identity_number, personal_code, first_name, last_name FROM users WHERE identity_number = $1
`

type FindUserByIdentityNumberRow struct {
	ID             uuid.UUID
	IdentityNumber string
	PersonalCode   string
	FirstName      string
	LastName       string
}

func (q *Queries) FindUserByIdentityNumber(ctx context.Context, identityNumber string) (FindUserByIdentityNumberRow, error) {
	row := q.db.QueryRow(ctx, findUserByIdentityNumber, identityNumber)
	var i FindUserByIdentityNumberRow
	err := row.Scan(
		&i.ID,
		&i.IdentityNumber,
		&i.PersonalCode,
		&i.FirstName,
		&i.LastName,
	)
	return i, err
}
