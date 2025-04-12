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

type scopesService struct {
	proto.UnimplementedScopeServiceServer
	scopes services.Scopes
	log    *logger.Logger
}

func NewScopes(scopes services.Scopes, log *logger.Logger) proto.ScopeServiceServer {
	return &scopesService{
		scopes: scopes,
		log:    log,
	}
}

//nolint:dupl
func (p *scopesService) List(ctx context.Context, req *proto.PaginatedListRequest) (*proto.ListScopesResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	pagination := &services.Pagination{
		Page:    req.Limit,
		PerPage: req.Offset,
	}

	rows, total, err := p.scopes.List(ctx, pagination)
	if err != nil {
		p.log.Error().Err(err).Msg("Failed to fetch scopes")

		switch {
		case errors.Is(err, errors.ErrFailedToFetchResults):
			return nil, status.Error(codes.Unavailable, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to fetch scopes")
		}
	}

	collection := make([]*proto.Scope, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, &proto.Scope{
			Id:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
		})
	}

	return &proto.ListScopesResponse{
		Data: collection,
		Meta: &proto.PaginationMeta{
			Page:  pagination.Page,
			Per:   pagination.PerPage,
			Total: total,
		},
	}, nil
}

func (p *scopesService) Get(ctx context.Context, req *proto.GetScopeRequest) (*proto.GetScopeResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to parse scope ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	scope, err := p.scopes.FindById(ctx, id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to get scope")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to get scope")
		}
	}

	return &proto.GetScopeResponse{
		Data: &proto.Scope{
			Id:          scope.ID.String(),
			Name:        scope.Name,
			Description: scope.Description,
		},
	}, nil
}

func (p *scopesService) Create(ctx context.Context, req *proto.CreateScopeRequest) (*proto.CreateScopeResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	scope, err := p.scopes.Create(ctx, &models.Scope{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		p.log.Error().Err(err).Str("name", req.Name).Msg("Failed to create scope")

		switch {
		case errors.Is(err, errors.ErrFailedToCreateRecord):
			return nil, status.Error(codes.Internal, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to create scope")
		}
	}

	return &proto.CreateScopeResponse{
		Data: &proto.Scope{
			Id:          scope.ID.String(),
			Name:        scope.Name,
			Description: scope.Description,
		},
	}, nil
}

func (p *scopesService) Update(ctx context.Context, req *proto.UpdateScopeRequest) (*proto.UpdateScopeResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Invalid UUID format")
		return nil, status.Error(codes.InvalidArgument, "invalid scope id format")
	}

	scope, err := p.scopes.Update(ctx, &models.Scope{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to update scope")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, errors.ErrFailedToUpdateRecord):
			return nil, status.Error(codes.Internal, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to update scope")
		}
	}

	return &proto.UpdateScopeResponse{
		Data: &proto.Scope{
			Id:          scope.ID.String(),
			Name:        scope.Name,
			Description: scope.Description,
		},
	}, nil
}

//nolint:dupl
func (p *scopesService) Delete(ctx context.Context, req *proto.DeleteScopeRequest) (*emptypb.Empty, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to parse scope ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = p.scopes.Delete(ctx, id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to delete scope")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to delete scope")
		}
	}

	return &emptypb.Empty{}, nil
}
