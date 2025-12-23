package main

import (
	"github.com/mark3labs/mcp-go/mcp/server"
)

func main() {
	// Create MCP Server
	server := server.NewMCPServer(
		"HelloMCPServer",
		"1.0.0",
		"An example MCP server in Go",
		server.WithToolCapabilities(false),
	)
}