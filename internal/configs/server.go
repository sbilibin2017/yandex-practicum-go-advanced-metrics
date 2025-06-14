package configs

// ServerConfig holds configuration settings for the server,
// including the network address to listen on and the logging level.
type ServerConfig struct {
	Address  string // Network address the server will listen on (e.g., ":8080")
	LogLevel string // Log verbosity level (e.g., "info", "debug", "error")
}

// ServerOption defines a functional option for configuring ServerConfig.
type ServerOption func(*ServerConfig)

// NewServerConfig creates a new ServerConfig instance applying the given options.
//
// Each ServerOption is applied in order to customize the configuration.
// If no options are provided, returns a ServerConfig with default zero-values.
func NewServerConfig(opts ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
