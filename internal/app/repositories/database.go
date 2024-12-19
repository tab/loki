package repositories

import (
	"context"

	"github.com/exaring/otelpgx"
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
	CreateUserTokens(ctx context.Context, params dto.CreateUserParams) (*models.User, error)
	CreateUserRole(ctx context.Context, params db.CreateUserRoleParams) error
	CreateUserScope(ctx context.Context, params db.CreateUserScopeParams) error

	FindUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindUserByIdentityNumber(ctx context.Context, identityNumber string) (*models.User, error)

	FindRoleByName(ctx context.Context, name string) (*models.Role, error)

	FindScopeByName(ctx context.Context, name string) (*models.Scope, error)

	FindUserRoles(ctx context.Context, id uuid.UUID) ([]models.Role, error)
	FindUserPermissions(ctx context.Context, id uuid.UUID) ([]models.Permission, error)
	FindUserScopes(ctx context.Context, id uuid.UUID) ([]models.Scope, error)
}

type database struct {
	db      *pgxpool.Pool
	queries *db.Queries
}

func NewDatabase(cfg *config.Config) (Database, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
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

	_, err = q.UpsertUserRoleByName(ctx, db.UpsertUserRoleByNameParams{
		UserID: user.ID,
		Name:   models.UserRoleType,
	})
	if err != nil {
		return nil, err
	}

	_, err = q.UpsertUserScopeByName(ctx, db.UpsertUserScopeByNameParams{
		UserID: user.ID,
		Name:   models.SelfServiceType,
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

func (d *database) CreateUserTokens(ctx context.Context, params dto.CreateUserParams) (*models.User, error) {
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

func (d *database) CreateUserRole(ctx context.Context, params db.CreateUserRoleParams) error {
	_, err := d.queries.CreateUserRole(ctx, params)
	return err
}

func (d *database) CreateUserScope(ctx context.Context, params db.CreateUserScopeParams) error {
	_, err := d.queries.CreateUserScope(ctx, params)
	return err
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

func (d *database) FindRoleByName(ctx context.Context, name string) (*models.Role, error) {
	record, err := d.queries.FindRoleByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &models.Role{
		ID:   record.ID,
		Name: record.Name,
	}, nil
}

func (d *database) FindScopeByName(ctx context.Context, name string) (*models.Scope, error) {
	record, err := d.queries.FindScopeByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &models.Scope{
		ID:   record.ID,
		Name: record.Name,
	}, nil
}

func (d *database) FindUserRoles(ctx context.Context, id uuid.UUID) ([]models.Role, error) {
	records, err := d.queries.FindUserRoles(ctx, id)
	if err != nil {
		return nil, err
	}

	roles := make([]models.Role, 0, len(records))
	for _, record := range records {
		roles = append(roles, models.Role{
			ID:   record.ID,
			Name: record.Name,
		})
	}

	return roles, nil
}

func (d *database) FindUserPermissions(ctx context.Context, id uuid.UUID) ([]models.Permission, error) {
	records, err := d.queries.FindUserPermissions(ctx, id)
	if err != nil {
		return nil, err
	}

	permissions := make([]models.Permission, 0, len(records))
	for _, record := range records {
		permissions = append(permissions, models.Permission{
			ID:   record.ID,
			Name: record.Name,
		})
	}

	return permissions, nil
}

func (d *database) FindUserScopes(ctx context.Context, id uuid.UUID) ([]models.Scope, error) {
	records, err := d.queries.FindUserScopes(ctx, id)
	if err != nil {
		return nil, err
	}

	scopes := make([]models.Scope, 0, len(records))
	for _, record := range records {
		scopes = append(scopes, models.Scope{
			ID:   record.ID,
			Name: record.Name,
		})
	}

	return scopes, nil
}
