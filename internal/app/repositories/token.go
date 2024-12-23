package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
)

type TokenRepository interface {
	Create(ctx context.Context, params db.CreateTokensParams) ([]models.Token, error)
}

type Token struct {
	client postgres.Postgres
}

func NewTokenRepository(client postgres.Postgres) TokenRepository {
	return &Token{client: client}
}

func (t *Token) Create(ctx context.Context, params db.CreateTokensParams) ([]models.Token, error) {
	tx, err := t.client.Db().Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := t.client.Queries().WithTx(tx)

	records, err := q.CreateTokens(ctx, db.CreateTokensParams{
		UserID:           params.UserID,
		AccessTokenValue: params.AccessTokenValue,
		AccessTokenExpiresAt: pgtype.Timestamp{
			Time:  time.Now().Add(models.AccessTokenExp),
			Valid: true,
		},
		RefreshTokenValue: params.RefreshTokenValue,
		RefreshTokenExpiresAt: pgtype.Timestamp{
			Time:  time.Now().Add(models.RefreshTokenExp),
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	tokens := make([]models.Token, 0, len(records))
	for _, record := range records {
		tokens = append(tokens, models.Token{
			ID:        record.ID,
			Type:      string(record.Type),
			Value:     record.Value,
			ExpiresAt: record.ExpiresAt.Time,
		})
	}

	return tokens, tx.Commit(ctx)
}
