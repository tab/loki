package services

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "go.uber.org/mock/gomock"

    "loki/internal/app/repositories"
)

func Test_HealthChecker_Ping(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    ctx := context.Background()
    repository := repositories.NewMockHealthRepository(ctrl)
    service := NewHealthChecker(repository)

    tests := []struct {
        name     string
        before   func()
        expected error
    }{
        {
            name: "Success",
            before: func() {
                repository.EXPECT().Ping(ctx).Return(nil)
            },
            expected: nil,
        },
        {
            name: "Error",
            before: func() {
                repository.EXPECT().Ping(ctx).Return(assert.AnError)
            },
            expected: assert.AnError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.before()

            result := service.Ping(ctx)
            assert.Equal(t, tt.expected, result)
        })
    }
}
