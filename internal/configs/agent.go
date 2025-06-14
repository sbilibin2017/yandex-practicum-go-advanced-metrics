package configs

// AgentConfig holds configuration settings for the metric agent.
type AgentConfig struct {
	// ServerAddress is the address of the metrics server.
	ServerAddress string

	// ServerEndpoint is the API endpoint for metric updates.
	ServerEndpoint string

	// LogLevel sets the logging level (e.g., "debug", "info").
	LogLevel string

	// PollInterval specifies the interval, in seconds, for polling metrics.
	PollInterval int

	// ReportInterval specifies the interval, in seconds, for reporting metrics.
	ReportInterval int

	// NumWorkers defines the number of concurrent workers for reporting.
	NumWorkers int
}

// AgentOption is a function type that modifies an AgentConfig.
// It allows optional configuration settings to be applied.
type AgentOption func(*AgentConfig)

// NewAgentConfig creates a new AgentConfig instance and applies
// any number of AgentOption functions to customize the configuration.
//
// Example:
//
//	cfg := NewAgentConfig(
//	    func(c *AgentConfig) { c.ServerAddress = "localhost:8080" },
//	    func(c *AgentConfig) { c.NumWorkers = 5 },
//	)
func NewAgentConfig(opts ...AgentOption) *AgentConfig {
	cfg := &AgentConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
