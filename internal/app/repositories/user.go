package repositories

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
)

type UserRepository interface {
	Create(ctx context.Context, params db.CreateUserParams) (*models.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByIdentityNumber(ctx context.Context, identityNumber string) (*models.User, error)
}

type user struct {
	client postgres.Postgres
}

func NewUserRepository(client postgres.Postgres) UserRepository {
	return &user{client: client}
}

func (u *user) Create(ctx context.Context, params db.CreateUserParams) (*models.User, error) {
	tx, err := u.client.Db().Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := u.client.Queries().WithTx(tx)

	result, err := q.CreateUser(ctx, db.CreateUserParams{
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
	})
	if err != nil {
		return nil, err
	}

	_, err = q.UpsertUserRoleByName(ctx, db.UpsertUserRoleByNameParams{
		UserID: result.ID,
		Name:   models.UserRoleType,
	})
	if err != nil {
		return nil, err
	}

	_, err = q.UpsertUserScopeByName(ctx, db.UpsertUserScopeByNameParams{
		UserID: result.ID,
		Name:   models.SelfServiceType,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:             result.ID,
		IdentityNumber: result.IdentityNumber,
		PersonalCode:   result.PersonalCode,
		FirstName:      result.FirstName,
		LastName:       result.LastName,
	}, tx.Commit(ctx)
}

func (u *user) FindById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	result, err := u.client.Queries().FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:             result.ID,
		IdentityNumber: result.IdentityNumber,
		PersonalCode:   result.PersonalCode,
		FirstName:      result.FirstName,
		LastName:       result.LastName,
	}, nil
}

func (u *user) FindByIdentityNumber(ctx context.Context, identityNumber string) (*models.User, error) {
	result, err := u.client.Queries().FindUserByIdentityNumber(ctx, identityNumber)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:             result.ID,
		IdentityNumber: result.IdentityNumber,
		PersonalCode:   result.PersonalCode,
		FirstName:      result.FirstName,
		LastName:       result.LastName,
	}, nil
}
