package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
)

// authKey is a custom context key for storing the auth token.
type authKey struct{}

// withAuthKey adds an auth key to the context.
func withAuthKey(ctx context.Context, auth string) context.Context {
	return context.WithValue(ctx, authKey{}, auth)
}

// authFromRequest extracts the auth token from the request headers.
func authFromRequest(ctx context.Context, r *http.Request) context.Context {
	return withAuthKey(ctx, r.Header.Get("Authorization"))
}

// authFromEnv extracts the auth token from the environment
func authFromEnv(ctx context.Context) context.Context {
	return withAuthKey(ctx, os.Getenv("API_KEY"))
}

// tokenFromContext extracts the auth token from the context.
// This can be used by tools to extract the token regardless of the
// transport being used by the server.
func tokenFromContext(ctx context.Context) (string, error) {
	auth, ok := ctx.Value(authKey{}).(string)
	if !ok {
		return "", fmt.Errorf("missing auth")
	}
	return auth, nil
}

type LinuxTerminalMCPServer struct {
	server *server.MCPServer
}

func NewLinuxTerminalMCPServer() *LinuxTerminalMCPServer {
	mcpServer := server.NewMCPServer(
		"linux-terminal",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)
	mcpServer.AddTool(
		mcp.NewTool("linux_terminal",
			mcp.WithDescription("在宿主机器上执行一个命令，命令将会直接在宿主机器上执行，请注意这个操作拥有很高权限，不要执行危险的命令，仅在需要时执行这个命令"),
			mcp.WithString("command",
				mcp.Description("The command to execute"),
				mcp.Required(),
			),
			mcp.WithString("dir",
				mcp.Description("The working directory to execute the command in, defaults to the current directory"),
			),
		),
		handleLinuxCommand)
	return &LinuxTerminalMCPServer{
		server: mcpServer,
	}
}

func handleLinuxCommand(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Handle the command execution here
	// For example, you can use exec.Command to run the command and capture the output
	command, ok := request.Params.Arguments["command"].(string)
	if !ok {
		return mcp.NewToolResultText("error, command can not be empty"), nil
	}
	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	dir, ok := request.Params.Arguments["dir"].(string)
	if ok {
		cmd.Dir = dir
	}
	slog.Info("Executing: %s", cmd)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	outChan := make(chan string)
	sb := bytes.NewBufferString("")
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			slog.Info(line)
			outChan <- line
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			slog.Info(line)
			outChan <- line
		}
	}()

	go func() {
		if err := cmd.Wait(); err != nil {
			slog.Error("Command execution failed", "error", err)
			outChan <- fmt.Sprintf("Command execution failed: %v", err)
		}
		close(outChan)
	}()

	for line := range outChan {
		sb.WriteString(line + "\n")
	}
	return mcp.NewToolResultText(sb.String()), nil
}

func (s *LinuxTerminalMCPServer) ServeSSE(addr string) *server.SSEServer {
	return server.NewSSEServer(s.server,
		server.WithBaseURL(fmt.Sprintf("http://%s", addr)),
		server.WithSSEContextFunc(authFromRequest),
	)
}

func (s *LinuxTerminalMCPServer) ServeStdio() error {
	return server.ServeStdio(s.server, server.WithStdioContextFunc(authFromEnv))
}

func main() {
	var transport string
	flag.StringVar(&transport, "t", "sse", "Transport type (stdio or sse)")
	flag.Parse()

	s := NewLinuxTerminalMCPServer()

	switch transport {
	case "stdio":
		if err := s.ServeStdio(); err != nil {
			slog.Error("Server error", "error", err)
		}
	case "sse":
		sseServer := s.ServeSSE("localhost:8088")
		slog.Info("SSE server listening on :8088")
		if err := sseServer.Start(":8088"); err != nil {
			slog.Error("Server error: %v", "error", err)
		}
	default:
		slog.Error("Invalid transport type: %s. Must be 'stdio' or 'sse'", "type", transport)
	}
}
