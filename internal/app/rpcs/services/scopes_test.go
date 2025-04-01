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

func Test_Scopes_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	scopes := services.NewMockScopes(ctrl)
	log := logger.NewLogger()
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
						ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
						Name:        "read:self",
						Description: "Read own data",
					},
					{
						ID:          uuid.MustParse("10000000-1000-1000-3000-000000000002"),
						Name:        "write:self",
						Description: "Write own data",
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
						Id:          "10000000-1000-1000-3000-000000000001",
						Name:        "read:self",
						Description: "Read own data",
					},
					{
						Id:          "10000000-1000-1000-3000-000000000002",
						Name:        "write:self",
						Description: "Write own data",
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
				scopes.EXPECT().List(ctx, gomock.Any()).Return(nil, uint64(0), errors.ErrFailedToFetchResults)
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

	ctx := context.Background()
	scopes := services.NewMockScopes(ctrl)
	log := logger.NewLogger()
	service := NewScopes(scopes, log)

	id := uuid.MustParse("10000000-1000-1000-3000-000000000001")

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
					Name:        "read:self",
					Description: "Read own data",
				}, nil)
			},
			req: &proto.GetScopeRequest{
				Id: id.String(),
			},
			expected: &proto.GetScopeResponse{
				Data: &proto.Scope{
					Id:          id.String(),
					Name:        "read:self",
					Description: "Read own data",
				},
			},
			error: false,
		},
		{
			name: "Not Found",
			before: func() {
				scopes.EXPECT().FindById(ctx, id).Return(nil, errors.ErrScopeNotFound)
			},
			req: &proto.GetScopeRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Invalid ID Format",
			req: &proto.GetScopeRequest{
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

			result, err := service.Get(ctx, tt.req)

			if tt.error {
				st, _ := status.FromError(err)
				assert.Equal(t, tt.code, st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Data.Id, result.Data.Id)
				assert.Equal(t, tt.expected.Data.Name, result.Data.Name)
				assert.Equal(t, tt.expected.Data.Description, result.Data.Description)
			}
		})
	}
}

func Test_Scopes_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	scopes := services.NewMockScopes(ctrl)
	log := logger.NewLogger()
	service := NewScopes(scopes, log)

	id := uuid.MustParse("10000000-1000-1000-3000-000000000001")

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
					Name:        "read:self",
					Description: "Read own data",
				}, nil)
			},
			req: &proto.CreateScopeRequest{
				Name:        "read:self",
				Description: "Read own data",
			},
			expected: &proto.CreateScopeResponse{
				Data: &proto.Scope{
					Id:          id.String(),
					Name:        "read:self",
					Description: "Read own data",
				},
			},
			error: false,
		},
		{
			name: "Internal Error",
			before: func() {
				scopes.EXPECT().Create(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			req: &proto.CreateScopeRequest{
				Name:        "read:self",
				Description: "Read own data",
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
		{
			name: "Validation Error",
			req: &proto.CreateScopeRequest{
				Name:        "",
				Description: "Read own data",
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

			result, err := service.Create(ctx, tt.req)

			if tt.error {
				st, _ := status.FromError(err)
				assert.Equal(t, tt.code, st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Data.Id, result.Data.Id)
				assert.Equal(t, tt.expected.Data.Name, result.Data.Name)
				assert.Equal(t, tt.expected.Data.Description, result.Data.Description)
			}
		})
	}
}

func Test_Scopes_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	scopes := services.NewMockScopes(ctrl)
	log := logger.NewLogger()
	service := NewScopes(scopes, log)

	id := uuid.MustParse("10000000-1000-1000-3000-000000000001")

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
					Name:        "read:self",
					Description: "Read own data updated",
				}, nil)
			},
			req: &proto.UpdateScopeRequest{
				Id:          id.String(),
				Name:        "read:self",
				Description: "Read own data updated",
			},
			expected: &proto.UpdateScopeResponse{
				Data: &proto.Scope{
					Id:          id.String(),
					Name:        "read:self",
					Description: "Read own data updated",
				},
			},
			error: false,
		},
		{
			name: "Not Found",
			before: func() {
				scopes.EXPECT().Update(ctx, gomock.Any()).Return(nil, errors.ErrScopeNotFound)
			},
			req: &proto.UpdateScopeRequest{
				Id:          id.String(),
				Name:        "read:self",
				Description: "Read own data updated",
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Internal Error",
			before: func() {
				scopes.EXPECT().Update(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			req: &proto.UpdateScopeRequest{
				Id:          id.String(),
				Name:        "read:self",
				Description: "Read own data updated",
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
		{
			name: "Invalid ID Format",
			req: &proto.UpdateScopeRequest{
				Id:          "invalid-uuid",
				Name:        "read:self",
				Description: "Read own data updated",
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

			result, err := service.Update(ctx, tt.req)

			if tt.error {
				st, _ := status.FromError(err)
				assert.Equal(t, tt.code, st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Data.Id, result.Data.Id)
				assert.Equal(t, tt.expected.Data.Name, result.Data.Name)
				assert.Equal(t, tt.expected.Data.Description, result.Data.Description)
			}
		})
	}
}

func Test_Scopes_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	scopes := services.NewMockScopes(ctrl)
	log := logger.NewLogger()
	service := NewScopes(scopes, log)

	id := uuid.MustParse("10000000-1000-1000-3000-000000000001")

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
			name: "Not Found",
			before: func() {
				scopes.EXPECT().Delete(ctx, id).Return(false, errors.ErrScopeNotFound)
			},
			req: &proto.DeleteScopeRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Internal Error",
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
		{
			name: "Invalid ID Format",
			req: &proto.DeleteScopeRequest{
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
