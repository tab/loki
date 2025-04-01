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

func Test_Permissions_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	permissions := services.NewMockPermissions(ctrl)
	log := logger.NewLogger()
	service := NewPermissions(permissions, log)

	tests := []struct {
		name     string
		before   func()
		request  *proto.PaginatedListRequest
		expected *proto.ListPermissionsResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				permissions.EXPECT().List(ctx, gomock.Any()).Return([]models.Permission{
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
			expected: &proto.ListPermissionsResponse{
				Data: []*proto.Permission{
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
				permissions.EXPECT().List(ctx, gomock.Any()).Return(nil, uint64(0), errors.ErrFailedToFetchResults)
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
				for i, permission := range tt.expected.Data {
					assert.Equal(t, permission.Id, result.Data[i].Id)
					assert.Equal(t, permission.Name, result.Data[i].Name)
					assert.Equal(t, permission.Description, result.Data[i].Description)
				}
			}
		})
	}
}

func Test_Permissions_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	permissions := services.NewMockPermissions(ctrl)
	log := logger.NewLogger()
	service := NewPermissions(permissions, log)

	id := uuid.MustParse("10000000-1000-1000-3000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.GetPermissionRequest
		expected *proto.GetPermissionResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				permissions.EXPECT().FindById(ctx, id).Return(&models.Permission{
					ID:          id,
					Name:        "read:self",
					Description: "Read own data",
				}, nil)
			},
			req: &proto.GetPermissionRequest{
				Id: id.String(),
			},
			expected: &proto.GetPermissionResponse{
				Data: &proto.Permission{
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
				permissions.EXPECT().FindById(ctx, id).Return(nil, errors.ErrPermissionNotFound)
			},
			req: &proto.GetPermissionRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Invalid ID Format",
			req: &proto.GetPermissionRequest{
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

func Test_Permissions_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	permissions := services.NewMockPermissions(ctrl)
	log := logger.NewLogger()
	service := NewPermissions(permissions, log)

	id := uuid.MustParse("10000000-1000-1000-3000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.CreatePermissionRequest
		expected *proto.CreatePermissionResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				permissions.EXPECT().Create(ctx, gomock.Any()).Return(&models.Permission{
					ID:          id,
					Name:        "read:self",
					Description: "Read own data",
				}, nil)
			},
			req: &proto.CreatePermissionRequest{
				Name:        "read:self",
				Description: "Read own data",
			},
			expected: &proto.CreatePermissionResponse{
				Data: &proto.Permission{
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
				permissions.EXPECT().Create(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			req: &proto.CreatePermissionRequest{
				Name:        "read:self",
				Description: "Read own data",
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
		{
			name: "Validation Error",
			req: &proto.CreatePermissionRequest{
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

func Test_Permissions_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	permissions := services.NewMockPermissions(ctrl)
	log := logger.NewLogger()
	service := NewPermissions(permissions, log)

	id := uuid.MustParse("10000000-1000-1000-3000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.UpdatePermissionRequest
		expected *proto.UpdatePermissionResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				permissions.EXPECT().Update(ctx, gomock.Any()).Return(&models.Permission{
					ID:          id,
					Name:        "read:self",
					Description: "Read own data updated",
				}, nil)
			},
			req: &proto.UpdatePermissionRequest{
				Id:          id.String(),
				Name:        "read:self",
				Description: "Read own data updated",
			},
			expected: &proto.UpdatePermissionResponse{
				Data: &proto.Permission{
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
				permissions.EXPECT().Update(ctx, gomock.Any()).Return(nil, errors.ErrPermissionNotFound)
			},
			req: &proto.UpdatePermissionRequest{
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
				permissions.EXPECT().Update(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			req: &proto.UpdatePermissionRequest{
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
			req: &proto.UpdatePermissionRequest{
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

func Test_Permissions_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	permissions := services.NewMockPermissions(ctrl)
	log := logger.NewLogger()
	service := NewPermissions(permissions, log)

	id := uuid.MustParse("10000000-1000-1000-3000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.DeletePermissionRequest
		expected *emptypb.Empty
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				permissions.EXPECT().Delete(ctx, id).Return(true, nil)
			},
			req: &proto.DeletePermissionRequest{
				Id: id.String(),
			},
			expected: &emptypb.Empty{},
			error:    false,
		},
		{
			name: "Not Found",
			before: func() {
				permissions.EXPECT().Delete(ctx, id).Return(false, errors.ErrPermissionNotFound)
			},
			req: &proto.DeletePermissionRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Internal Error",
			before: func() {
				permissions.EXPECT().Delete(ctx, id).Return(false, assert.AnError)
			},
			req: &proto.DeletePermissionRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
		{
			name: "Invalid ID Format",
			req: &proto.DeletePermissionRequest{
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
