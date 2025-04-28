package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/graydovee/xiaoshi/pkg/config"
	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ChatBot struct {
	ctx    context.Context
	cancel context.CancelFunc

	openaiClient openai.Client
	mcpClients   map[string]*mcpclient.Client
	model        string

	tools         []openai.ChatCompletionToolParam
	toolsToClient map[string]*mcpclient.Client
}

func (c *ChatBot) Model() string {
	return c.model
}

func (c *ChatBot) SetModel(model string) {
	c.model = model
}

func NewChatBot(config *config.MCPConfig) (*ChatBot, error) {
	chatBot := &ChatBot{
		model: config.LLM.Model,
	}

	chatBot.openaiClient = openai.NewClient(
		option.WithAPIKey(config.LLM.ApiKey),
		option.WithBaseURL(config.LLM.BaseURL),
	)

	clients, err := NewMCPClients(config)
	if err != nil {
		return nil, err
	}
	chatBot.mcpClients = clients

	return chatBot, nil
}
func (c *ChatBot) Initialize() error {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.toolsToClient = make(map[string]*mcpclient.Client)

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "xiaoxi",
		Version: "1.0.0",
	}

	for _, cli := range c.mcpClients {
		if err := cli.Start(c.ctx); err != nil {
			return err
		}
		initResult, err := cli.Initialize(c.ctx, initRequest)
		if err != nil {
			return fmt.Errorf("failed to initialize: %w", err)
		}
		slog.Info("Initialized with server",
			"name", initResult.ServerInfo.Name,
			"version", initResult.ServerInfo.Version,
		)

		// get tools
		toolsRequest := mcp.ListToolsRequest{}
		tools, err := cli.ListTools(c.ctx, toolsRequest)
		if err != nil {
			return fmt.Errorf("failed to list tools: %w", err)
		}
		for _, tool := range tools.Tools {
			slog.Info("Tool found", "name", tool.Name, "description", tool.Description)
		}

		for _, tool := range tools.Tools {
			c.tools = append(c.tools, openai.ChatCompletionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        tool.Name,
					Description: openai.String(tool.Description),
					Parameters: openai.FunctionParameters{
						"type":       tool.InputSchema.Type,
						"properties": tool.InputSchema.Properties,
						"required":   tool.InputSchema.Required,
					},
				},
			})
			c.toolsToClient[tool.Name] = cli
		}
	}
	return nil
}

func (c *ChatBot) Completion(ctx context.Context, prompt string, history History) (chan string, error) {
	history.AddHistory(openai.UserMessage(prompt))

	msgChan := make(chan string, 10)
	go func() {
		defer close(msgChan)
		for {
			toolsResp, err := c.completion(ctx, msgChan, history)
			if err != nil {
				slog.Error("Failed to get completion", "error", err)
			}
			if !toolsResp {
				break
			}
		}
	}()
	return msgChan, nil
}

func (c *ChatBot) completion(ctx context.Context, msgChan chan string, history History) (bool, error) {
	chatCompletion, err := c.openaiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: history.GetHistory(),
		Tools:    c.tools,
		Model:    c.model,
	})
	if err != nil {
		return false, err
	}

	if len(chatCompletion.Choices) == 0 {
		return false, fmt.Errorf("no choices returned")
	}
	resp := chatCompletion.Choices[0]
	if resp.Message.Content != "" {
		history.AddHistory(openai.AssistantMessage(resp.Message.Content))
		msgChan <- resp.Message.Content
	}

	if len(resp.Message.ToolCalls) > 0 {
		history.AddHistory(resp.Message.ToParam())
	}
	var toolsResp bool
	for _, fc := range resp.Message.ToolCalls {
		fn := fc.Function.Name
		cli, ok := c.toolsToClient[fn]
		if !ok {
			slog.Error("Failed to find client", "function", fn)
			continue
		}

		slog.Info("Calling tool", "name", fn, "arguments", fc.Function.Arguments)

		mcpRequest := mcp.CallToolRequest{}
		mcpRequest.Params.Name = fc.Function.Name
		err = json.Unmarshal([]byte(fc.Function.Arguments), &mcpRequest.Params.Arguments)
		if err != nil {
			slog.Error("Failed to parse arguments", "function", fn, "error", err)
			continue
		}
		result, err := cli.CallTool(c.ctx, mcpRequest)
		if err != nil {
			slog.Error("Failed to call tool", "function", fn, "error", err)
			continue
		}
		if result.IsError {
			slog.Error("Tool returned error", "function", fn, "meta", result.Meta)
			continue
		}

		for _, content := range result.Content {
			switch tc := content.(type) {
			case mcp.TextContent:
				history.AddHistory(openai.ToolMessage(tc.Text, fc.ID))
				toolsResp = true
			case mcp.EmbeddedResource:
				jsonBytes, _ := json.MarshalIndent(content, "", "  ")
				msgChan <- string(jsonBytes)
			default:
				slog.Warn("Unknown content type", "type", tc)
				jsonBytes, _ := json.MarshalIndent(content, "", "  ")
				slog.Info("Content", "data", string(jsonBytes))
			}
		}
	}
	return toolsResp, nil
}
