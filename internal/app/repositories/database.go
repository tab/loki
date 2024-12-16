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
	CreateUser(ctx context.Context, params db.CreateUserTokensParams) (*models.User, error)
	CreateUserTokens(ctx context.Context, params dto.CreateUserTokensParams) (*models.User, error)
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

func (d *database) CreateUser(ctx context.Context, params db.CreateUserTokensParams) (*models.User, error) {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := d.queries.WithTx(tx)

	user, err := q.CreateUser(ctx, db.CreateUserTokensParams{
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
	})
	if err != nil {
		return nil, err
	}

	role, err := q.FindRoleByName(ctx, models.UserRoleType)
	if err != nil {
		return nil, err
	}

	_, err = q.CreateUserRole(ctx, db.CreateUserRoleParams{
		UserID: user.ID,
		RoleID: role.ID,
	})
	if err != nil {
		return nil, err
	}

	scope, err := q.FindScopeByName(ctx, models.SelfServiceType)
	if err != nil {
		return nil, err
	}

	_, err = q.CreateUserScope(ctx, db.CreateUserScopeParams{
		UserID:  user.ID,
		ScopeID: scope.ID,
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
	}, tx.Commit(ctx)
}

func (d *database) CreateUserTokens(ctx context.Context, params dto.CreateUserTokensParams) (*models.User, error) {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := d.queries.WithTx(tx)

	user, err := q.FindUserByIdentityNumber(ctx, params.IdentityNumber)
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
