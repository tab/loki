package services

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"loki/internal/app/errors"
	proto "loki/internal/app/rpcs/proto/sso/v1"
	"loki/internal/app/services"
	"loki/pkg/logger"
)

type tokensService struct {
	proto.UnimplementedTokenServiceServer
	tokens services.Tokens
	log    *logger.Logger
}

func NewTokens(tokens services.Tokens, log *logger.Logger) proto.TokenServiceServer {
	return &tokensService{
		tokens: tokens,
		log:    log,
	}
}

//nolint:dupl
func (p *tokensService) List(ctx context.Context, req *proto.PaginatedListRequest) (*proto.ListTokensResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	pagination := &services.Pagination{
		Page:    req.Limit,
		PerPage: req.Offset,
	}

	rows, total, err := p.tokens.List(ctx, pagination)
	if err != nil {
		p.log.Error().Err(err).Msg("Failed to fetch tokens")

		switch {
		case errors.Is(err, errors.ErrFailedToFetchResults):
			return nil, status.Error(codes.Unavailable, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to fetch tokens")
		}
	}

	collection := make([]*proto.Token, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, &proto.Token{
			Id:        row.ID.String(),
			UserId:    row.UserId.String(),
			Type:      row.Type,
			Value:     row.Value,
			ExpiresAt: timestamppb.New(row.ExpiresAt),
		})
	}

	return &proto.ListTokensResponse{
		Data: collection,
		Meta: &proto.PaginationMeta{
			Page:  pagination.Page,
			Per:   pagination.PerPage,
			Total: total,
		},
	}, nil
}

//nolint:dupl
func (p *tokensService) Delete(ctx context.Context, req *proto.DeleteTokenRequest) (*emptypb.Empty, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidArguments.Error())
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to parse token ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = p.tokens.Delete(ctx, id)
	if err != nil {
		p.log.Error().Err(err).Str("id", req.Id).Msg("Failed to delete token")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to delete token")
		}
	}

	return &emptypb.Empty{}, nil
}
