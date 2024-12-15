// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: token.sql

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
