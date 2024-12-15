package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/app/repositories/db"
	"loki/internal/config"
)

type Database interface {
	CreateUser(ctx context.Context, params db.CreateUserParams) (*models.User, error)
	CreateOrUpdateUserWithTokens(ctx context.Context, params dto.CreateUserParams) (*models.User, error)
	RefreshUserTokens(ctx context.Context, params dto.RefreshTokenParams) (*models.User, error)
	FindUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindUserByIdentityNumber(ctx context.Context, identityNumber string) (*models.User, error)
}

type database struct {
	db      *pgxpool.Pool
	queries *db.Queries
}

func NewDatabase(cfg *config.Config) (Database, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	queries := db.New(pool)

	return &database{
		db:      pool,
		queries: queries,
	}, nil
}

func (d *database) CreateUser(ctx context.Context, params db.CreateUserParams) (*models.User, error) {
	record, err := d.queries.CreateUser(ctx, db.CreateUserParams{
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
	})

	return &models.User{
		ID:             record.ID,
		IdentityNumber: record.IdentityNumber,
		PersonalCode:   record.PersonalCode,
		FirstName:      record.FirstName,
		LastName:       record.LastName,
	}, err
}

func (d *database) CreateOrUpdateUserWithTokens(ctx context.Context, params dto.CreateUserParams) (*models.User, error) {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := d.queries.WithTx(tx)

	user, err := q.CreateUser(ctx, db.CreateUserParams{
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := q.CreateToken(ctx, db.CreateTokenParams{
		UserID: user.ID,
		Type:   db.TokenType(params.AccessToken.Type),
		Value:  params.AccessToken.Value,
		ExpiresAt: pgtype.Timestamp{
			Time:  params.AccessToken.ExpiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := q.CreateToken(ctx, db.CreateTokenParams{
		UserID: user.ID,
		Type:   db.TokenType(params.RefreshToken.Type),
		Value:  params.RefreshToken.Value,
		ExpiresAt: pgtype.Timestamp{
			Time:  params.RefreshToken.ExpiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:             user.ID,
		IdentityNumber: user.IdentityNumber,
		PersonalCode:   user.PersonalCode,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		AccessToken:    accessToken.Value,
		RefreshToken:   refreshToken.Value,
	}, tx.Commit(ctx)
}

func (d *database) RefreshUserTokens(ctx context.Context, params dto.RefreshTokenParams) (*models.User, error) {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := d.queries.WithTx(tx)

	user, err := q.FindUserById(ctx, params.UserId)
	if err != nil {
		return nil, err
	}

	accessToken, err := q.CreateToken(ctx, db.CreateTokenParams{
		UserID: user.ID,
		Type:   db.TokenType(params.AccessToken.Type),
		Value:  params.AccessToken.Value,
		ExpiresAt: pgtype.Timestamp{
			Time:  params.AccessToken.ExpiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := q.CreateToken(ctx, db.CreateTokenParams{
		UserID: user.ID,
		Type:   db.TokenType(params.RefreshToken.Type),
		Value:  params.RefreshToken.Value,
		ExpiresAt: pgtype.Timestamp{
			Time:  params.RefreshToken.ExpiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:             user.ID,
		IdentityNumber: user.IdentityNumber,
		PersonalCode:   user.PersonalCode,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		AccessToken:    accessToken.Value,
		RefreshToken:   refreshToken.Value,
	}, tx.Commit(ctx)
}

func (d *database) FindUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	record, err := d.queries.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:             record.ID,
		IdentityNumber: record.IdentityNumber,
		PersonalCode:   record.PersonalCode,
		FirstName:      record.FirstName,
		LastName:       record.LastName,
	}, nil
}

func (d *database) FindUserByIdentityNumber(ctx context.Context, identityNumber string) (*models.User, error) {
	record, err := d.queries.FindUserByIdentityNumber(ctx, identityNumber)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:             record.ID,
		IdentityNumber: record.IdentityNumber,
		PersonalCode:   record.PersonalCode,
		FirstName:      record.FirstName,
		LastName:       record.LastName,
	}, nil
}
