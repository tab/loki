package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"loki/internal/config"
)

func Test_NewRedisClient(t *testing.T) {
	type args struct {
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
				cfg: &config.Config{
					RedisURI: "redis://localhost:6379",
				},
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewRedisClient(tt.args.cfg)
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}
