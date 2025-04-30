package config

type Config struct {
	Memory  MemoryConfig  `json:"memory" yaml:"memory"`
	OneBot  OneBotConfig  `json:"oneBot" yaml:"oneBot"`
	MCP     MCPConfig     `json:"mcp" yaml:"mcp"`
	Healthz HealthzConfig `json:"healthz" yaml:"healthz"`
}

type HealthzConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Addr    string `json:"addr" yaml:"addr"`
	Port    int    `json:"port" yaml:"port"`
	Pattern string `json:"pattern" yaml:"pattern"`
}

type MemoryConfig struct {
	MessageLimit  int `json:"messageLimit" yaml:"messageLimit"`
	ExpireSeconds int `json:"expireSeconds" yaml:"expireSeconds"`
}

type OneBotConfig struct {
	Id    int64        `json:"id" yaml:"id"`
	Ws    *WebSocket   `json:"ws,omitempty" yaml:"ws,omitempty"`
	Limit *LimitConfig `json:"limit,omitempty" yaml:"limit,omitempty"`
}

type LimitConfig struct {
	Enabled   bool    `json:"enabled" yaml:"enabled"`
	Frequency float64 `json:"frequency" yaml:"frequency"`
	Bucket    int     `json:"bucket" yaml:"bucket"`
}

type WebSocket struct {
	Addr  string `json:"addr" yaml:"addr"`
	Port  int    `json:"port" yaml:"port"`
	Token string `json:"token" yaml:"token"`
}

type MCPConfig struct {
	SystemPrompt string               `json:"systemPrompt" yaml:"systemPrompt"`
	LLM          LLM                  `json:"llm" yaml:"llm"`
	McpServers   map[string]MCPServer `json:"mcpServers" yaml:"mcpServers"`
}

type MCPServer struct {
	Command string            `json:"command" yaml:"command"`
	Args    []string          `json:"args" yaml:"args"`
	Env     map[string]string `json:"env" yaml:"env"`

	Url     string            `json:"url" yaml:"url"`
	Headers map[string]string `json:"headers" yaml:"headers"`

	Disabled bool `json:"disabled" yaml:"disabled"`
}

type LLM struct {
	BaseURL string `json:"baseUrl" yaml:"baseUrl"`
	ApiKey  string `json:"apiKey" yaml:"apiKey"`
	Model   string `json:"model" yaml:"model"`
}
