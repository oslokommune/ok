package cmd

import (
	"github.com/spf13/cobra"
)

// mcpCommand represents the mcp command
var mcpCommand = &cobra.Command{
	Use:   "mcp",
	Short: "Model Context Protocol (MCP) server operations",
	Long: `The mcp command provides functionality for running Model Context Protocol servers.
	
MCP enables LLM applications to securely connect to and interact with external data sources 
and tools through a standardized protocol.`,
}