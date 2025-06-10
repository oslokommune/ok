package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

// ServeCommand represents the serve command
var ServeCommand = &cobra.Command{
	Use:   "serve",
	Short: "Start an MCP server",
	Long: `Start a Model Context Protocol server that provides tools and resources to LLM applications.

The server runs over stdio by default, which is the standard transport for MCP servers.`,
	RunE: runServe,
}

func runServe(cmd *cobra.Command, args []string) error {
	// Create a new MCP server
	s := server.NewMCPServer(
		"OK MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// Add hello world tool
	helloTool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	// Add tool handler
	s.AddTool(helloTool, helloHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}