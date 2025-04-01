package repositories

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
)

type UserRepository interface {
	List(ctx context.Context, limit, offset uint64) ([]models.User, uint64, error)
	Create(ctx context.Context, params db.CreateUserParams) (*models.User, error)
	Update(ctx context.Context, params db.UpdateUserParams) (*models.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)

	FindByIdentityNumber(ctx context.Context, identityNumber string) (*models.User, error)
	FindUserDetailsById(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type user struct {
	client postgres.Postgres
}

func NewUserRepository(client postgres.Postgres) UserRepository {
	return &user{client: client}
}

func (u *user) List(ctx context.Context, limit, offset uint64) ([]models.User, uint64, error) {
	rows, err := u.client.Queries().FindUsers(ctx, db.FindUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	users := make([]models.User, 0, len(rows))
	var total uint64

	if len(rows) > 0 {
		total = rows[0].Total
	}

	for _, row := range rows {
		users = append(users, models.User{
			ID:             row.ID,
			IdentityNumber: row.IdentityNumber.String,
			PersonalCode:   row.PersonalCode.String,
			FirstName:      row.FirstName.String,
			LastName:       row.LastName.String,
		})
	}

	return users, total, err
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

func (u *user) Update(ctx context.Context, params db.UpdateUserParams) (*models.User, error) {
	tx, err := u.client.Db().Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := u.client.Queries().WithTx(tx)

	result, err := q.UpdateUser(ctx, db.UpdateUserParams{
		ID:             params.ID,
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
	})
	if err != nil {
		return nil, err
	}

	_, err = q.CreateUserRoles(ctx, db.CreateUserRolesParams{
		UserID:  result.ID,
		RoleIds: params.RoleIDs,
	})
	if err != nil {
		return nil, err
	}

	_, err = q.CreateUserScopes(ctx, db.CreateUserScopesParams{
		UserID:   result.ID,
		ScopeIds: params.ScopeIDs,
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

func (u *user) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	err := u.client.Queries().DeleteUser(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
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

func (u *user) FindUserDetailsById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	result, err := u.client.Queries().FindUserDetailsById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:             result.ID,
		IdentityNumber: result.IdentityNumber,
		PersonalCode:   result.PersonalCode,
		FirstName:      result.FirstName,
		LastName:       result.LastName,
		RoleIDs:        result.RoleIds,
		ScopeIDs:       result.ScopeIds,
	}, nil
}
