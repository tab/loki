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

func Test_Users_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	users := services.NewMockUsers(ctrl)
	log := logger.NewLogger()
	service := NewUsers(users, log)

	tests := []struct {
		name     string
		before   func()
		request  *proto.PaginatedListRequest
		expected *proto.ListUsersResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				users.EXPECT().List(ctx, gomock.Any()).Return([]models.User{
					{
						ID:             uuid.MustParse("10000000-1000-1000-1234-000000000001"),
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
					},
					{
						ID:             uuid.MustParse("10000000-1000-1000-1234-000000000002"),
						IdentityNumber: "PNOEE-987654321",
						PersonalCode:   "987654321",
						FirstName:      "Jane",
						LastName:       "TESTNUMBER",
					},
				}, uint64(2), nil)
			},
			request: &proto.PaginatedListRequest{
				Limit:  1,
				Offset: 10,
			},
			expected: &proto.ListUsersResponse{
				Data: []*proto.User{
					{
						Id:             "10000000-1000-1000-1234-000000000001",
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
					},
					{
						Id:             "10000000-1000-1000-1234-000000000002",
						IdentityNumber: "PNOEE-987654321",
						PersonalCode:   "987654321",
						FirstName:      "Jane",
						LastName:       "TESTNUMBER",
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
				users.EXPECT().List(ctx, gomock.Any()).Times(0)
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
				users.EXPECT().List(ctx, gomock.Any()).Return(nil, uint64(0), errors.ErrFailedToFetchResults)
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
				users.EXPECT().List(ctx, gomock.Any()).Return(nil, uint64(0), assert.AnError)
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
				for i, user := range tt.expected.Data {
					assert.Equal(t, user.Id, result.Data[i].Id)
					assert.Equal(t, user.IdentityNumber, result.Data[i].IdentityNumber)
					assert.Equal(t, user.PersonalCode, result.Data[i].PersonalCode)
					assert.Equal(t, user.FirstName, result.Data[i].FirstName)
					assert.Equal(t, user.LastName, result.Data[i].LastName)
				}
			}
		})
	}
}

func Test_Users_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	users := services.NewMockUsers(ctrl)
	log := logger.NewLogger()
	service := NewUsers(users, log)

	id := uuid.MustParse("10000000-1000-1000-1234-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.GetUserRequest
		expected *proto.GetUserResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				users.EXPECT().FindById(ctx, id).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}, nil)
			},
			req: &proto.GetUserRequest{
				Id: id.String(),
			},
			expected: &proto.GetUserResponse{
				Data: &proto.User{
					Id:             id.String(),
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				},
			},
			error: false,
		},
		{
			name:   "Invalid ID format",
			before: func() {},
			req: &proto.GetUserRequest{
				Id: "invalid-uuid",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Not Found",
			before: func() {
				users.EXPECT().FindById(ctx, id).Return(nil, errors.ErrRecordNotFound)
			},
			req: &proto.GetUserRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Error",
			before: func() {
				users.EXPECT().FindById(ctx, id).Return(nil, assert.AnError)
			},
			req: &proto.GetUserRequest{
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
				assert.Equal(t, tt.expected.Data.Id, result.Data.Id)
				assert.Equal(t, tt.expected.Data.IdentityNumber, result.Data.IdentityNumber)
				assert.Equal(t, tt.expected.Data.PersonalCode, result.Data.PersonalCode)
				assert.Equal(t, tt.expected.Data.FirstName, result.Data.FirstName)
				assert.Equal(t, tt.expected.Data.LastName, result.Data.LastName)
			}
		})
	}
}

func Test_Users_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	users := services.NewMockUsers(ctrl)
	log := logger.NewLogger()
	service := NewUsers(users, log)

	id := uuid.MustParse("10000000-1000-1000-1234-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.CreateUserRequest
		expected *proto.CreateUserResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				users.EXPECT().Create(gomock.Any(), &models.User{
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}, nil)
			},
			req: &proto.CreateUserRequest{
				IdentityNumber: "PNOEE-60001017869",
				PersonalCode:   "60001017869",
				FirstName:      "EID2016",
				LastName:       "TESTNUMBER",
			},
			expected: &proto.CreateUserResponse{
				Data: &proto.User{
					Id:             id.String(),
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				},
			},
			error: false,
		},
		{
			name:   "Validation error",
			before: func() {},
			req: &proto.CreateUserRequest{
				IdentityNumber: "",
				PersonalCode:   "60001017869",
				FirstName:      "EID2016",
				LastName:       "TESTNUMBER",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Error",
			before: func() {
				users.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.ErrFailedToCreateRecord)
			},
			req: &proto.CreateUserRequest{
				IdentityNumber: "PNOEE-60001017869",
				PersonalCode:   "60001017869",
				FirstName:      "EID2016",
				LastName:       "TESTNUMBER",
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
		{
			name: "Internal error",
			before: func() {
				users.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
			req: &proto.CreateUserRequest{
				IdentityNumber: "PNOEE-60001017869",
				PersonalCode:   "60001017869",
				FirstName:      "EID2016",
				LastName:       "TESTNUMBER",
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
				assert.Equal(t, tt.expected.Data.Id, result.Data.Id)
				assert.Equal(t, tt.expected.Data.IdentityNumber, result.Data.IdentityNumber)
				assert.Equal(t, tt.expected.Data.PersonalCode, result.Data.PersonalCode)
				assert.Equal(t, tt.expected.Data.FirstName, result.Data.FirstName)
				assert.Equal(t, tt.expected.Data.LastName, result.Data.LastName)
			}
		})
	}
}

func Test_Users_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	users := services.NewMockUsers(ctrl)
	log := logger.NewLogger()
	service := NewUsers(users, log)

	id := uuid.MustParse("10000000-1000-1000-1234-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.UpdateUserRequest
		expected *proto.UpdateUserResponse
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				users.EXPECT().Update(ctx, gomock.Any()).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "JOHN",
					LastName:       "DOE",
				}, nil)
			},
			req: &proto.UpdateUserRequest{
				Id:             id.String(),
				IdentityNumber: "PNOEE-60001017869",
				PersonalCode:   "60001017869",
				FirstName:      "JOHN",
				LastName:       "DOE",
			},
			expected: &proto.UpdateUserResponse{
				Data: &proto.User{
					Id:             id.String(),
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "JOHN",
					LastName:       "DOE",
				},
			},
			error: false,
		},
		{
			name:   "Invalid ID format",
			before: func() {},
			req: &proto.UpdateUserRequest{
				Id:             "invalid-uuid",
				IdentityNumber: "PNOEE-60001017869",
				PersonalCode:   "60001017869",
				FirstName:      "JOHN",
				LastName:       "DOE",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Not Found",
			before: func() {
				users.EXPECT().Update(ctx, gomock.Any()).Return(nil, errors.ErrRecordNotFound)
			},
			req: &proto.UpdateUserRequest{
				Id:             id.String(),
				IdentityNumber: "PNOEE-60001017869",
				PersonalCode:   "60001017869",
				FirstName:      "JOHN",
				LastName:       "DOE",
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Error",
			before: func() {
				users.EXPECT().Update(ctx, gomock.Any()).Return(nil, errors.ErrFailedToUpdateRecord)
			},
			req: &proto.UpdateUserRequest{
				Id:             id.String(),
				IdentityNumber: "PNOEE-60001017869",
				PersonalCode:   "60001017869",
				FirstName:      "JOHN",
				LastName:       "DOE",
			},
			expected: nil,
			code:     codes.Internal,
			error:    true,
		},
		{
			name: "Internal error",
			before: func() {
				users.EXPECT().Update(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			req: &proto.UpdateUserRequest{
				Id:             id.String(),
				IdentityNumber: "PNOEE-60001017869",
				PersonalCode:   "60001017869",
				FirstName:      "JOHN",
				LastName:       "DOE",
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
				assert.Equal(t, tt.expected.Data.Id, result.Data.Id)
				assert.Equal(t, tt.expected.Data.IdentityNumber, result.Data.IdentityNumber)
				assert.Equal(t, tt.expected.Data.PersonalCode, result.Data.PersonalCode)
				assert.Equal(t, tt.expected.Data.FirstName, result.Data.FirstName)
				assert.Equal(t, tt.expected.Data.LastName, result.Data.LastName)
			}
		})
	}
}

func Test_Users_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	users := services.NewMockUsers(ctrl)
	log := logger.NewLogger()
	service := NewUsers(users, log)

	id := uuid.MustParse("10000000-1000-1000-1234-000000000001")

	tests := []struct {
		name     string
		before   func()
		req      *proto.DeleteUserRequest
		expected *emptypb.Empty
		code     codes.Code
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				users.EXPECT().Delete(ctx, id).Return(true, nil)
			},
			req: &proto.DeleteUserRequest{
				Id: id.String(),
			},
			expected: &emptypb.Empty{},
			error:    false,
		},
		{
			name:   "Invalid ID format",
			before: func() {},
			req: &proto.DeleteUserRequest{
				Id: "invalid-uuid",
			},
			expected: nil,
			code:     codes.InvalidArgument,
			error:    true,
		},
		{
			name: "Not Found",
			before: func() {
				users.EXPECT().Delete(ctx, id).Return(false, errors.ErrRecordNotFound)
			},
			req: &proto.DeleteUserRequest{
				Id: id.String(),
			},
			expected: nil,
			code:     codes.NotFound,
			error:    true,
		},
		{
			name: "Internal error",
			before: func() {
				users.EXPECT().Delete(ctx, id).Return(false, assert.AnError)
			},
			req: &proto.DeleteUserRequest{
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
