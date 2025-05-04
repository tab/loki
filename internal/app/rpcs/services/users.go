package services

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/app/services"
	"loki/internal/config/logger"
	proto "loki/internal/app/rpcs/proto/sso/v1"
)

type usersService struct {
	proto.UnimplementedUserServiceServer
	users services.Users
	log   *logger.Logger
}

func NewUsers(users services.Users, log *logger.Logger) proto.UserServiceServer {
	return &usersService{
		users: users,
		log:   log,
	}
}

//nolint:dupl
func (p *usersService) List(ctx context.Context, req *proto.PaginatedListRequest) (*proto.ListUsersResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	pagination := &services.Pagination{
		Page:    req.Limit,
		PerPage: req.Offset,
	}

	rows, total, err := p.users.List(ctx, pagination)
	if err != nil {
		p.log.Error().Err(err).Msg("Failed to fetch users")

		switch {
		case errors.Is(err, errors.ErrFailedToFetchResults):
			return nil, status.Error(codes.Unavailable, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to fetch users")
		}
	}

	collection := make([]*proto.User, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, &proto.User{
			Id:             row.ID.String(),
			IdentityNumber: row.IdentityNumber,
			PersonalCode:   row.PersonalCode,
			FirstName:      row.FirstName,
			LastName:       row.LastName,
		})
	}

	return &proto.ListUsersResponse{
		Data: collection,
		Meta: &proto.PaginationMeta{
			Page:  pagination.Page,
			Per:   pagination.PerPage,
			Total: total,
		},
	}, nil
}

func (p *usersService) Get(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to parse user ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := p.users.FindUserDetailsById(ctx, id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to get user")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to get user")
		}
	}

	roleIds := make([]string, 0, len(user.RoleIDs))
	for _, roleId := range user.RoleIDs {
		roleIds = append(roleIds, roleId.String())
	}

	scopeIds := make([]string, 0, len(user.ScopeIDs))
	for _, scopeId := range user.ScopeIDs {
		scopeIds = append(scopeIds, scopeId.String())
	}

	return &proto.GetUserResponse{
		Data: &proto.User{
			Id:             user.ID.String(),
			IdentityNumber: user.IdentityNumber,
			PersonalCode:   user.PersonalCode,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			RoleIds:        roleIds,
			ScopeIds:       scopeIds,
		},
	}, nil
}

func (p *usersService) Create(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	user, err := p.users.Create(ctx, &models.User{
		IdentityNumber: req.IdentityNumber,
		PersonalCode:   req.PersonalCode,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
	})
	if err != nil {
		p.log.Error().Err(err).Str("identity_number", req.IdentityNumber).Msg("Failed to create user")

		switch {
		case errors.Is(err, errors.ErrFailedToCreateRecord):
			return nil, status.Error(codes.Internal, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to create user")
		}
	}

	return &proto.CreateUserResponse{
		Data: &proto.User{
			Id:             user.ID.String(),
			IdentityNumber: user.IdentityNumber,
			PersonalCode:   user.PersonalCode,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
		},
	}, nil
}

func (p *usersService) Update(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Invalid UUID format")
		return nil, status.Error(codes.InvalidArgument, "invalid user id format")
	}

	roleIds := make([]uuid.UUID, 0, len(req.RoleIds))
	for _, roleId := range req.RoleIds {
		id, err := uuid.Parse(roleId)
		if err != nil {
			p.log.Error().Err(err).Str("id", roleId).Msg("Invalid UUID format")
			return nil, status.Error(codes.InvalidArgument, "invalid role id format")
		}
		roleIds = append(roleIds, id)
	}

	scopeIds := make([]uuid.UUID, 0, len(req.ScopeIds))
	for _, scopeId := range req.ScopeIds {
		id, err := uuid.Parse(scopeId)
		if err != nil {
			p.log.Error().Err(err).Str("id", scopeId).Msg("Invalid UUID format")
			return nil, status.Error(codes.InvalidArgument, "invalid scope id format")
		}
		scopeIds = append(scopeIds, id)
	}

	user, err := p.users.Update(ctx, &models.User{
		ID:             id,
		IdentityNumber: req.IdentityNumber,
		PersonalCode:   req.PersonalCode,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		RoleIDs:        roleIds,
		ScopeIDs:       scopeIds,
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to update user")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, errors.ErrFailedToUpdateRecord):
			return nil, status.Error(codes.Internal, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to update user")
		}
	}

	return &proto.UpdateUserResponse{
		Data: &proto.User{
			Id:             user.ID.String(),
			IdentityNumber: user.IdentityNumber,
			PersonalCode:   user.PersonalCode,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
		},
	}, nil
}

//nolint:dupl
func (p *usersService) Delete(ctx context.Context, req *proto.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to parse user ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = p.users.Delete(ctx, id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to delete user")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to delete user")
		}
	}

	return &emptypb.Empty{}, nil
}
