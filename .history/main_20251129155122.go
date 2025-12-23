package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP Server
	server := server.NewMCPServer(
		"HelloMCPServer",
		"1.0.0",
		"An example MCP server in Go",
		server.WithToolCapabilities(false),
	)

	// Add the Tool
	tool := mcp.NewTool("hello_tool",
		mcp.WithDescription("A tool that says hello"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name to greet")),
	)
}