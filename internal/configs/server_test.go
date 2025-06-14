package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig_Default(t *testing.T) {
	cfg := NewServerConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, "", cfg.Address) // По умолчанию пустая строка
}

func TestNewServerConfig_WithAddress(t *testing.T) {
	address := "localhost:8080"
	opt := func(cfg *ServerConfig) {
		cfg.Address = address
	}

	cfg := NewServerConfig(opt)
	assert.NotNil(t, cfg)
	assert.Equal(t, address, cfg.Address)
}
