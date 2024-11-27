package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/config"
	"loki/internal/config/server"
	"loki/pkg/logger"
)

func Test_NewApplication(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		expected Application
	}{
		{
			name:     "Success",
			expected: Application{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := NewApplication(ctx)
			assert.NoError(t, err)

			assert.NotNil(t, app)
			assert.NotNil(t, app.cfg)
			assert.NotNil(t, app.log)
			assert.NotNil(t, app.server)
		})
	}
}

func Test_Application_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := server.NewMockServer(ctrl)
	appLogger := logger.NewLogger()

	tests := []struct {
		name     string
		before   func()
		expected error
	}{
		{
			name: "Success",
			before: func() {
				mockServer.EXPECT().Run().DoAndReturn(func() error {
					time.Sleep(100 * time.Millisecond)
					return nil
				})
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			app := &Application{
				cfg:    &config.Config{},
				log:    appLogger,
				server: mockServer,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()

			err := app.Run(ctx)

			if tt.expected != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
