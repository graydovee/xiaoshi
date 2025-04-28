package mcp

import (
	"github.com/graydovee/xiaoshi/pkg/config"
	mcpclient "github.com/mark3labs/mcp-go/client"
)

func NewMCPClients(config *config.MCPConfig) (map[string]*mcpclient.Client, error) {
	res := make(map[string]*mcpclient.Client)
	for name, cfg := range config.McpServers {
		if cfg.Disabled {
			continue
		}

		client, err := NewMcpClient(cfg)
		if err != nil {
			return nil, err
		}
		res[name] = client
	}
	return res, nil
}

func NewMcpClient(server config.MCPServer) (*mcpclient.Client, error) {
	if server.Url == "" {
		env := make([]string, 0, len(server.Env))
		for k, v := range server.Env {
			env = append(env, k+"="+v)
		}
		return mcpclient.NewStdioMCPClient(
			server.Command,
			env,
			server.Args...,
		)
	} else {
		return mcpclient.NewSSEMCPClient(server.Url, mcpclient.WithHeaders(server.Header))
	}
}

func ReadCompletionResponse(ch <-chan string) string {
	sb := ""
	for line := range ch {
		sb += line
	}
	return sb
}
