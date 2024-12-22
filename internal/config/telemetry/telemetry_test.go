package telemetry

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"loki/internal/config"
)

func Test_NewTelemetry(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *config.Config
	}

	tests := []struct {
		name string
		args args
		err  bool
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				cfg: &config.Config{
					AppName:      "loki",
					TelemetryURI: "http://localhost:4317",
				},
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewTelemetry(tt.args.ctx, tt.args.cfg)
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}
