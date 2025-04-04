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
	proto "loki/internal/app/rpcs/proto/sso/v1"
	"loki/internal/app/services"
	"loki/pkg/logger"
)

type rolesService struct {
	proto.UnimplementedRoleServiceServer
	roles services.Roles
	log   *logger.Logger
}

func NewRoles(roles services.Roles, log *logger.Logger) proto.RoleServiceServer {
	return &rolesService{
		roles: roles,
		log:   log,
	}
}

//nolint:dupl
func (p *rolesService) List(ctx context.Context, req *proto.PaginatedListRequest) (*proto.ListRolesResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	pagination := &services.Pagination{
		Page:    req.Limit,
		PerPage: req.Offset,
	}

	rows, total, err := p.roles.List(ctx, pagination)
	if err != nil {
		p.log.Error().Err(err).Msg("Failed to fetch roles")

		switch {
		case errors.Is(err, errors.ErrFailedToFetchResults):
			return nil, status.Error(codes.Unavailable, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to fetch roles")
		}
	}

	collection := make([]*proto.Role, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, &proto.Role{
			Id:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
		})
	}

	return &proto.ListRolesResponse{
		Data: collection,
		Meta: &proto.PaginationMeta{
			Page:  pagination.Page,
			Per:   pagination.PerPage,
			Total: total,
		},
	}, nil
}

func (p *rolesService) Get(ctx context.Context, req *proto.GetRoleRequest) (*proto.GetRoleResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to parse role ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	role, err := p.roles.FindById(ctx, id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to get role")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to get role")
		}
	}

	return &proto.GetRoleResponse{
		Data: &proto.Role{
			Id:            role.ID.String(),
			Name:          role.Name,
			Description:   role.Description,
			PermissionIds: []string{},
		},
	}, nil
}

func (p *rolesService) Create(ctx context.Context, req *proto.CreateRoleRequest) (*proto.CreateRoleResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	permissionIDs := make([]uuid.UUID, 0, len(req.PermissionIds))
	for _, permissionId := range req.PermissionIds {
		id, err := uuid.Parse(permissionId)
		if err != nil {
			p.log.Error().Err(err).Str("permission_id", permissionId).Msg("Invalid permission ID format")
			return nil, status.Error(codes.InvalidArgument, "invalid permission ID format")
		}
		permissionIDs = append(permissionIDs, id)
	}

	role, err := p.roles.Create(ctx, &models.Role{
		Name:          req.Name,
		Description:   req.Description,
		PermissionIDs: permissionIDs,
	})
	if err != nil {
		p.log.Error().Err(err).Str("name", req.Name).Msg("Failed to create role")

		switch {
		case errors.Is(err, errors.ErrFailedToCreateRecord):
			return nil, status.Error(codes.Internal, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to create role")
		}
	}

	return &proto.CreateRoleResponse{
		Data: &proto.Role{
			Id:          role.ID.String(),
			Name:        role.Name,
			Description: role.Description,
		},
	}, nil
}

func (p *rolesService) Update(ctx context.Context, req *proto.UpdateRoleRequest) (*proto.UpdateRoleResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	permissionIDs := make([]uuid.UUID, 0, len(req.PermissionIds))
	for _, permissionId := range req.PermissionIds {
		id, err := uuid.Parse(permissionId)
		if err != nil {
			p.log.Error().Err(err).Str("permission_id", permissionId).Msg("Invalid permission ID format")
			return nil, status.Error(codes.InvalidArgument, "invalid permission ID format")
		}
		permissionIDs = append(permissionIDs, id)
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Invalid UUID format")
		return nil, status.Error(codes.InvalidArgument, "invalid role id format")
	}

	role, err := p.roles.Update(ctx, &models.Role{
		ID:            id,
		Name:          req.Name,
		Description:   req.Description,
		PermissionIDs: permissionIDs,
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to update role")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, errors.ErrFailedToUpdateRecord):
			return nil, status.Error(codes.Internal, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to update role")
		}
	}

	return &proto.UpdateRoleResponse{
		Data: &proto.Role{
			Id:          role.ID.String(),
			Name:        role.Name,
			Description: role.Description,
		},
	}, nil
}

//nolint:dupl
func (p *rolesService) Delete(ctx context.Context, req *proto.DeleteRoleRequest) (*emptypb.Empty, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to parse role ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = p.roles.Delete(ctx, id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to delete role")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to delete role")
		}
	}

	return &emptypb.Empty{}, nil
}
