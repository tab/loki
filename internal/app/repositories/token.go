package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
)

type TokenRepository interface {
	List(ctx context.Context, limit, offset int32) ([]models.Token, int, error)
	Create(ctx context.Context, params db.CreateTokensParams) ([]models.Token, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Token, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}

type token struct {
	client postgres.Postgres
}

func NewTokenRepository(client postgres.Postgres) TokenRepository {
	return &token{client: client}
}

func (t *token) List(ctx context.Context, limit, offset int32) ([]models.Token, int, error) {
	rows, err := t.client.Queries().FindTokens(ctx, db.FindTokensParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	tokens := make([]models.Token, 0, len(rows))
	var total int

	if len(rows) > 0 {
		total = int(rows[0].Total)
	}

	for _, row := range rows {
		tokens = append(tokens, models.Token{
			ID:        row.ID,
			UserId:    row.UserID,
			Type:      string(row.Type),
			Value:     row.Value,
			ExpiresAt: row.ExpiresAt.Time,
		})
	}

	return tokens, total, err
}

func (t *token) Create(ctx context.Context, params db.CreateTokensParams) ([]models.Token, error) {
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

func (t *token) FindById(ctx context.Context, id uuid.UUID) (*models.Token, error) {
	result, err := t.client.Queries().FindTokenById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.Token{
		ID:        result.ID,
		UserId:    result.UserID,
		Type:      string(result.Type),
		Value:     result.Value,
		ExpiresAt: result.ExpiresAt.Time,
	}, nil
}

func (t *token) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	err := t.client.Queries().DeleteToken(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}
