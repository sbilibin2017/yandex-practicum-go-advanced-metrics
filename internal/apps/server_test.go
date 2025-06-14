package apps

import (
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestNewServerApp(t *testing.T) {
	tests := []struct {
		name    string
		config  *configs.ServerConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &configs.ServerConfig{
				Address: ":8080",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := NewServerApp(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, app)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, app)
				assert.Equal(t, tt.config.Address, app.Addr)
			}
		})
	}
}
