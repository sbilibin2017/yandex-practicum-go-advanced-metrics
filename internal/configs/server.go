package configs

// ServerConfig holds configuration settings for the server,
// including the network address to listen on and the logging level.
type ServerConfig struct {
	// Address is the network address the server listens on (e.g., ":8080").
	Address string

	// LogLevel sets the log verbosity level (e.g., "info", "debug", "error").
	LogLevel string
}

// ServerOption defines a functional option for configuring ServerConfig.
// This enables flexible and optional configuration by applying a series
// of functions that modify the ServerConfig instance.
type ServerOption func(*ServerConfig)

// NewServerConfig creates a new ServerConfig and applies any number
// of ServerOption functions to customize its fields.
//
// If no options are provided, the returned ServerConfig will have zero values.
//
// Example:
//
//	cfg := NewServerConfig(
//	    func(c *ServerConfig) { c.Address = ":8080" },
//	    func(c *ServerConfig) { c.LogLevel = "debug" },
//	)
func NewServerConfig(opts ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
