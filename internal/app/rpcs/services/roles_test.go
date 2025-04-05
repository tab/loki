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

func Test_Roles_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	roles := services.NewMockRoles(ctrl)
	log := logger.NewLogger()
	service := NewRoles(roles, log)

	tests := []struct {
		name     string
		before   func()
		request  *proto.PaginatedListRequest
		expected *proto.ListRolesResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				roles.EXPECT().List(ctx, gomock.Any()).Return([]models.Role{
					{
						ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
						Name:        models.AdminRoleType,
						Description: "Admin role",
					},
					{
						ID:          uuid.MustParse("10000000-1000-1000-1000-000000000002"),
						Name:        models.ManagerRoleType,
						Description: "Manager role",
					},
					{
						ID:          uuid.MustParse("10000000-1000-1000-1000-000000000003"),
						Name:        models.UserRoleType,
						Description: "User role",
					},
				}, uint64(3), nil)
			},
			request: &proto.PaginatedListRequest{
				Limit:  1,
				Offset: 10,
			},
			expected: &proto.ListRolesResponse{
				Data: []*proto.Role{
					{
						Id:          "10000000-1000-1000-1000-000000000001",
						Name:        "admin",
						Description: "Admin role",
					},
					{
						Id:          "10000000-1000-1000-1000-000000000002",
						Name:        "manager",
						Description: "Manager role",
					},
					{
						Id:          "10000000-1000-1000-1000-000000000003",
						Name:        "user",
						Description: "User role",
					},
				},
				Meta: &proto.PaginationMeta{
					Page:  1,
					Per:   10,
					Total: 3,
				},
			},
			error: false,
		},
		{
			name: "Invalid Request",
			before: func() {
				roles.EXPECT().List(ctx, gomock.Any()).Times(0)
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
				roles.EXPECT().List(ctx, gomock.Any()).Return(nil, uint64(0), errors.ErrFailedToFetchResults)
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
				roles.EXPECT().List(ctx, gomock.Any()).Return(nil, uint64(0), assert.AnError)
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
				for i, role := range tt.expected.Data {
					assert.Equal(t, role.Id, result.Data[i].Id)
					assert.Equal(t, role.Name, result.Data[i].Name)
					assert.Equal(t, role.Description, result.Data[i].Description)
				}
			}
		})
	}
}

func Test_Roles_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	roles := services.NewMockRoles(ctrl)
	log := logger.NewLogger()
	service := NewRoles(roles, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")
	permissionIds := []uuid.UUID{
		uuid.MustParse("10000000-1000-1000-3000-000000000001"),
		uuid.MustParse("10000000-1000-1000-3000-000000000002"),
	}

	tests := []struct {
		name     string
		before   func()
		req      *proto.GetRoleRequest
		expected *proto.GetRoleResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				roles.EXPECT().FindRoleDetailsById(ctx, id).Return(&models.Role{
					ID:            id,
					Name:          models.AdminRoleType,
					Description:   "Admin role",
					PermissionIDs: permissionIds,
				}, nil)
			},
			req: &proto.GetRoleRequest{
				Id: id.String(),
			},
			expected: &proto.GetRoleResponse{
				Data: &proto.Role{
					Id:            id.String(),
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{"10000000-1000-1000-3000-000000000001", "10000000-1000-1000-3000-000000000002"},
				},
			},
			error: false,
		},
		{
			name: "Not Found",
			before: func() {
				roles.EXPECT().FindRoleDetailsById(ctx, id).Return(nil, errors.ErrRecordNotFound)
			},
			req: &proto.GetRoleRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name:   "Invalid ID format",
			before: func() {},
			req: &proto.GetRoleRequest{
				Id: "invalid-uuid",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Error",
			before: func() {
				roles.EXPECT().FindRoleDetailsById(ctx, id).Return(nil, assert.AnError)
			},
			req: &proto.GetRoleRequest{
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

func Test_Roles_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	roles := services.NewMockRoles(ctrl)
	log := logger.NewLogger()
	service := NewRoles(roles, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.CreateRoleRequest
		expected *proto.CreateRoleResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				roles.EXPECT().Create(ctx, gomock.Any()).Return(&models.Role{
					ID:          id,
					Name:        models.AdminRoleType,
					Description: "Admin role",
				}, nil)
			},
			req: &proto.CreateRoleRequest{
				Name:        "admin",
				Description: "Admin role",
			},
			expected: &proto.CreateRoleResponse{
				Data: &proto.Role{
					Id:          id.String(),
					Name:        "admin",
					Description: "Admin role",
				},
			},
			error: false,
		},
		{
			name:   "Validation error",
			before: func() {},
			req: &proto.CreateRoleRequest{
				Name:        "",
				Description: "Admin role",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Error",
			before: func() {
				roles.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.ErrFailedToCreateRecord)
			},
			req: &proto.CreateRoleRequest{
				Name:        "admin",
				Description: "Admin role",
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
		{
			name: "Internal error",
			before: func() {
				roles.EXPECT().Create(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			req: &proto.CreateRoleRequest{
				Name:        "admin",
				Description: "Admin role",
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

func Test_Roles_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	roles := services.NewMockRoles(ctrl)
	log := logger.NewLogger()
	service := NewRoles(roles, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.UpdateRoleRequest
		expected *proto.UpdateRoleResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				roles.EXPECT().Update(ctx, gomock.Any()).Return(&models.Role{
					ID:          id,
					Name:        models.AdminRoleType,
					Description: "Admin role updated",
				}, nil)
			},
			req: &proto.UpdateRoleRequest{
				Id:          id.String(),
				Name:        "admin",
				Description: "Admin role updated",
			},
			expected: &proto.UpdateRoleResponse{
				Data: &proto.Role{
					Id:          id.String(),
					Name:        "admin",
					Description: "Admin role updated",
				},
			},
			error: false,
		},
		{
			name:   "Validation error",
			before: func() {},
			req: &proto.UpdateRoleRequest{
				Id:   id.String(),
				Name: "",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name:   "Invalid ID format",
			before: func() {},
			req: &proto.UpdateRoleRequest{
				Id:          "invalid-uuid",
				Name:        "admin",
				Description: "Admin role updated",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Not Found",
			before: func() {
				roles.EXPECT().Update(ctx, gomock.Any()).Return(nil, errors.ErrRecordNotFound)
			},
			req: &proto.UpdateRoleRequest{
				Id:          id.String(),
				Name:        "admin",
				Description: "Admin role updated",
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Error",
			before: func() {
				roles.EXPECT().Update(ctx, gomock.Any()).Return(nil, errors.ErrFailedToUpdateRecord)
			},
			req: &proto.UpdateRoleRequest{
				Id:          id.String(),
				Name:        "admin",
				Description: "Admin role updated",
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
		{
			name: "Internal error",
			before: func() {
				roles.EXPECT().Update(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			req: &proto.UpdateRoleRequest{
				Id:          id.String(),
				Name:        "admin",
				Description: "Admin role updated",
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

func Test_Roles_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	roles := services.NewMockRoles(ctrl)
	log := logger.NewLogger()
	service := NewRoles(roles, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.DeleteRoleRequest
		expected *emptypb.Empty
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				roles.EXPECT().Delete(ctx, id).Return(true, nil)
			},
			req: &proto.DeleteRoleRequest{
				Id: id.String(),
			},
			expected: &emptypb.Empty{},
			error:    false,
		},
		{
			name:   "Invalid ID format",
			before: func() {},
			req: &proto.DeleteRoleRequest{
				Id: "invalid-uuid",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Not Found",
			before: func() {
				roles.EXPECT().Delete(ctx, id).Return(false, errors.ErrRecordNotFound)
			},
			req: &proto.DeleteRoleRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Internal error",
			before: func() {
				roles.EXPECT().Delete(ctx, id).Return(false, assert.AnError)
			},
			req: &proto.DeleteRoleRequest{
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
