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
	"loki/internal/config"
	"loki/internal/config/logger"
)

func Test_Scopes_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	scopes := services.NewMockScopes(ctrl)
	service := NewScopes(scopes, log)

	tests := []struct {
		name     string
		before   func()
		request  *proto.PaginatedListRequest
		expected *proto.ListScopesResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				scopes.EXPECT().List(ctx, gomock.Any()).Return([]models.Scope{
					{
						ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
						Name:        models.SsoServiceType,
						Description: "SSO-service scope",
					},
					{
						ID:          uuid.MustParse("10000000-1000-1000-2000-000000000002"),
						Name:        models.SelfServiceType,
						Description: "Self-service scope",
					},
				}, uint64(2), nil)
			},
			request: &proto.PaginatedListRequest{
				Limit:  1,
				Offset: 10,
			},
			expected: &proto.ListScopesResponse{
				Data: []*proto.Scope{
					{
						Id:          "10000000-1000-1000-2000-000000000001",
						Name:        "sso-service",
						Description: "SSO-service scope",
					},
					{
						Id:          "10000000-1000-1000-2000-000000000002",
						Name:        "self-service",
						Description: "Self-service scope",
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
			name: "Invalid Request",
			before: func() {
				scopes.EXPECT().List(ctx, gomock.Any()).Times(0)
			},
			request: &proto.PaginatedListRequest{
				Limit:  0,
				Offset: 10,
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Failed to fetch results",
			before: func() {
				scopes.EXPECT().List(ctx, gomock.Any()).Return(nil, uint64(0), errors.ErrFailedToFetchResults)
			},
			request: &proto.PaginatedListRequest{
				Limit:  1,
				Offset: 10,
			},
			expected: nil,
			code:     codes.Unavailable,
			error:    true,
		},
		{
			name: "Error",
			before: func() {
				scopes.EXPECT().List(ctx, gomock.Any()).Return(nil, uint64(0), assert.AnError)
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
				for i, scope := range tt.expected.Data {
					assert.Equal(t, scope.Id, result.Data[i].Id)
					assert.Equal(t, scope.Name, result.Data[i].Name)
					assert.Equal(t, scope.Description, result.Data[i].Description)
				}
			}
		})
	}
}

func Test_Scopes_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	scopes := services.NewMockScopes(ctrl)
	service := NewScopes(scopes, log)

	id := uuid.MustParse("10000000-1000-1000-2000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.GetScopeRequest
		expected *proto.GetScopeResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				scopes.EXPECT().FindById(ctx, id).Return(&models.Scope{
					ID:          id,
					Name:        models.SsoServiceType,
					Description: "SSO-service scope",
				}, nil)
			},
			req: &proto.GetScopeRequest{
				Id: id.String(),
			},
			expected: &proto.GetScopeResponse{
				Data: &proto.Scope{
					Id:          id.String(),
					Name:        "sso-service",
					Description: "SSO-service scope",
				},
			},
			error: false,
		},
		{
			name:   "Invalid ID format",
			before: func() {},
			req: &proto.GetScopeRequest{
				Id: "invalid-uuid",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Not Found",
			before: func() {
				scopes.EXPECT().FindById(ctx, id).Return(nil, errors.ErrRecordNotFound)
			},
			req: &proto.GetScopeRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Error",
			before: func() {
				scopes.EXPECT().FindById(ctx, id).Return(nil, assert.AnError)
			},
			req: &proto.GetScopeRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Get(ctx, tt.req)

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

func Test_Scopes_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	scopes := services.NewMockScopes(ctrl)
	service := NewScopes(scopes, log)

	id := uuid.MustParse("10000000-1000-1000-2000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.CreateScopeRequest
		expected *proto.CreateScopeResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				scopes.EXPECT().Create(ctx, gomock.Any()).Return(&models.Scope{
					ID:          id,
					Name:        models.SsoServiceType,
					Description: "SSO-service scope",
				}, nil)
			},
			req: &proto.CreateScopeRequest{
				Name:        "sso-service",
				Description: "SSO-service scope",
			},
			expected: &proto.CreateScopeResponse{
				Data: &proto.Scope{
					Id:          id.String(),
					Name:        "sso-service",
					Description: "SSO-service scope",
				},
			},
			error: false,
		},
		{
			name:   "Validation error",
			before: func() {},
			req: &proto.CreateScopeRequest{
				Name:        "",
				Description: "SSO-service scope",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Error",
			before: func() {
				scopes.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.ErrFailedToCreateRecord)
			},
			req: &proto.CreateScopeRequest{
				Name:        "sso-service",
				Description: "SSO-service scope",
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
		{
			name: "Internal error",
			before: func() {
				scopes.EXPECT().Create(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			req: &proto.CreateScopeRequest{
				Name:        "sso-service",
				Description: "SSO-service scope",
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Create(ctx, tt.req)

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

func Test_Scopes_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	scopes := services.NewMockScopes(ctrl)
	service := NewScopes(scopes, log)

	id := uuid.MustParse("10000000-1000-1000-2000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.UpdateScopeRequest
		expected *proto.UpdateScopeResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				scopes.EXPECT().Update(ctx, gomock.Any()).Return(&models.Scope{
					ID:          id,
					Name:        models.SsoServiceType,
					Description: "SSO-service scope updated",
				}, nil)
			},
			req: &proto.UpdateScopeRequest{
				Id:          id.String(),
				Name:        "sso-service",
				Description: "SSO-service scope updated",
			},
			expected: &proto.UpdateScopeResponse{
				Data: &proto.Scope{
					Id:          id.String(),
					Name:        "sso-service",
					Description: "SSO-service scope updated",
				},
			},
			error: false,
		},
		{
			name:   "Invalid ID format",
			before: func() {},
			req: &proto.UpdateScopeRequest{
				Id:          "invalid-uuid",
				Name:        "sso-service",
				Description: "SSO-service scope updated",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Not Found",
			before: func() {
				scopes.EXPECT().Update(ctx, gomock.Any()).Return(nil, errors.ErrRecordNotFound)
			},
			req: &proto.UpdateScopeRequest{
				Id:          id.String(),
				Name:        "sso-service",
				Description: "SSO-service scope updated",
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Error",
			before: func() {
				scopes.EXPECT().Update(ctx, gomock.Any()).Return(nil, errors.ErrFailedToUpdateRecord)
			},
			req: &proto.UpdateScopeRequest{
				Id:          id.String(),
				Name:        "sso-service",
				Description: "SSO-service scope updated",
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
		{
			name: "Internal error",
			before: func() {
				scopes.EXPECT().Update(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			req: &proto.UpdateScopeRequest{
				Id:          id.String(),
				Name:        "sso-service",
				Description: "SSO-service scope updated",
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Update(ctx, tt.req)

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

func Test_Scopes_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	scopes := services.NewMockScopes(ctrl)
	service := NewScopes(scopes, log)

	id := uuid.MustParse("10000000-1000-1000-2000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.DeleteScopeRequest
		expected *emptypb.Empty
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				scopes.EXPECT().Delete(ctx, id).Return(true, nil)
			},
			req: &proto.DeleteScopeRequest{
				Id: id.String(),
			},
			expected: &emptypb.Empty{},
			error:    false,
		},
		{
			name:   "Invalid ID format",
			before: func() {},
			req: &proto.DeleteScopeRequest{
				Id: "invalid-uuid",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Not Found",
			before: func() {
				scopes.EXPECT().Delete(ctx, id).Return(false, errors.ErrRecordNotFound)
			},
			req: &proto.DeleteScopeRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Internal error",
			before: func() {
				scopes.EXPECT().Delete(ctx, id).Return(false, assert.AnError)
			},
			req: &proto.DeleteScopeRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

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
