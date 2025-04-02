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

type permissionsService struct {
	proto.UnimplementedPermissionServiceServer
	permissions services.Permissions
	log         *logger.Logger
}

func NewPermissions(permissions services.Permissions, log *logger.Logger) proto.PermissionServiceServer {
	return &permissionsService{
		permissions: permissions,
		log:         log,
	}
}

//nolint:dupl
func (p *permissionsService) List(ctx context.Context, req *proto.PaginatedListRequest) (*proto.ListPermissionsResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrFailedToFetchResults.Error())
	}

	pagination := &services.Pagination{
		Page:    req.Limit,
		PerPage: req.Offset,
	}

	rows, total, err := p.permissions.List(ctx, pagination)
	if err != nil {
		p.log.Error().Err(err).Msg("Failed to fetch permissions")
		return nil, status.Error(codes.Internal, "failed to fetch permissions")
	}

	collection := make([]*proto.Permission, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, &proto.Permission{
			Id:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
		})
	}

	return &proto.ListPermissionsResponse{
		Data: collection,
		Meta: &proto.PaginationMeta{
			Page:  pagination.Page,
			Per:   pagination.PerPage,
			Total: total,
		},
	}, nil
}

func (p *permissionsService) Get(ctx context.Context, req *proto.GetPermissionRequest) (*proto.GetPermissionResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidAttributes.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to parse permission ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	permission, err := p.permissions.FindById(ctx, id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to get permission")
		if errors.Is(err, errors.ErrPermissionNotFound) {
			return nil, status.Error(codes.NotFound, "permission not found")
		}
		return nil, status.Error(codes.Internal, "failed to get permission")
	}

	return &proto.GetPermissionResponse{
		Data: &proto.Permission{
			Id:          permission.ID.String(),
			Name:        permission.Name,
			Description: permission.Description,
		},
	}, nil
}

func (p *permissionsService) Create(ctx context.Context, req *proto.CreatePermissionRequest) (*proto.CreatePermissionResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidAttributes.Error())
	}

	permission, err := p.permissions.Create(ctx, &models.Permission{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		p.log.Error().Err(err).Str("name", req.Name).Msg("Failed to create permission")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.CreatePermissionResponse{
		Data: &proto.Permission{
			Id:          permission.ID.String(),
			Name:        permission.Name,
			Description: permission.Description,
		},
	}, nil
}

func (p *permissionsService) Update(ctx context.Context, req *proto.UpdatePermissionRequest) (*proto.UpdatePermissionResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidAttributes.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Invalid UUID format")
		return nil, status.Error(codes.InvalidArgument, "invalid permission id format")
	}

	permission, err := p.permissions.Update(ctx, &models.Permission{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to update permission")

		switch {
		case errors.Is(err, errors.ErrPermissionNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to update permission")
		}
	}

	return &proto.UpdatePermissionResponse{
		Data: &proto.Permission{
			Id:          permission.ID.String(),
			Name:        permission.Name,
			Description: permission.Description,
		},
	}, nil
}

//nolint:dupl
func (p *permissionsService) Delete(ctx context.Context, req *proto.DeletePermissionRequest) (*emptypb.Empty, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidAttributes.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to parse permission ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = p.permissions.Delete(ctx, id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to delete permission")

		switch {
		case errors.Is(err, errors.ErrPermissionNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to delete permission")
		}
	}

	return &emptypb.Empty{}, nil
}
