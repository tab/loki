package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	proto "loki/internal/app/rpcs/proto/sso/v1"
	"loki/internal/app/services"
	"loki/pkg/logger"
)

func Test_Tokens_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	tokens := services.NewMockTokens(ctrl)
	log := logger.NewLogger()
	service := NewTokens(tokens, log)

	tests := []struct {
		name     string
		before   func()
		request  *proto.PaginatedListRequest
		expected *proto.ListTokensResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				tokens.EXPECT().List(ctx, gomock.Any()).Return([]models.Token{
					{
						ID:     uuid.MustParse("10000000-1000-1000-6000-000000000001"),
						UserId: uuid.MustParse("10000000-1000-1000-1234-000000000001"),
						Type:   models.AccessTokenType,
						Value:  "access-token-value",
					},
					{
						ID:     uuid.MustParse("10000000-1000-1000-6000-000000000002"),
						UserId: uuid.MustParse("10000000-1000-1000-1234-000000000002"),
						Type:   models.RefreshTokenType,
						Value:  "refresh-token-value",
					},
				}, uint64(2), nil)
			},
			request: &proto.PaginatedListRequest{
				Limit:  1,
				Offset: 10,
			},
			expected: &proto.ListTokensResponse{
				Data: []*proto.Token{
					{
						Id:     "10000000-1000-1000-6000-000000000001",
						UserId: "10000000-1000-1000-1234-000000000001",
						Type:   "access_token",
						Value:  "access-token-value",
					},
					{
						Id:     "10000000-1000-1000-6000-000000000002",
						UserId: "10000000-1000-1000-1234-000000000002",
						Type:   "refresh_token",
						Value:  "refresh-token-value",
					},
				},
				Meta: &proto.PaginationMeta{
					Page:  1,
					Per:   10,
					Total: 2,
				},
			},
			error: false,
		},
		{
			name: "Error",
			before: func() {
				tokens.EXPECT().List(ctx, gomock.Any()).Return(nil, uint64(0), errors.ErrFailedToFetchResults)
			},
			request: &proto.PaginatedListRequest{
				Limit:  1,
				Offset: 10,
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.List(ctx, tt.request)

			if tt.error {
				st, _ := status.FromError(err)
				assert.Equal(t, tt.code, st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expected.Data), len(result.Data))
				assert.Equal(t, tt.expected.Meta.Total, result.Meta.Total)
				for i, token := range tt.expected.Data {
					assert.Equal(t, token.Id, result.Data[i].Id)
					assert.Equal(t, token.UserId, result.Data[i].UserId)
					assert.Equal(t, token.Type, result.Data[i].Type)
					assert.Equal(t, token.Value, result.Data[i].Value)
				}
			}
		})
	}
}

func Test_Tokens_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	tokens := services.NewMockTokens(ctrl)
	log := logger.NewLogger()
	service := NewTokens(tokens, log)

	id := uuid.MustParse("10000000-1000-1000-6000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.DeleteTokenRequest
		expected *emptypb.Empty
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				tokens.EXPECT().Delete(ctx, id).Return(true, nil)
			},
			req: &proto.DeleteTokenRequest{
				Id: id.String(),
			},
			expected: &emptypb.Empty{},
			error:    false,
		},
		{
			name: "Not Found",
			before: func() {
				tokens.EXPECT().Delete(ctx, id).Return(false, errors.ErrTokenNotFound)
			},
			req: &proto.DeleteTokenRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Internal Error",
			before: func() {
				tokens.EXPECT().Delete(ctx, id).Return(false, assert.AnError)
			},
			req: &proto.DeleteTokenRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
		{
			name: "Invalid ID Format",
			req: &proto.DeleteTokenRequest{
				Id: "invalid-uuid",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}

			result, err := service.Delete(ctx, tt.req)

			if tt.error {
				st, _ := status.FromError(err)
				assert.Equal(t, tt.code, st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
