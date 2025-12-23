package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP Server
	srv := server.NewMCPServer(
		"HelloMCPServer",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add the Tool
	tool := mcp.NewTool("hello_tool",
		mcp.WithDescription("A tool that says hello"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name to greet")),
	)

	// Tool Handler
	srv.AddTool(tool, helloHandler)

	// Start the stdio server
	if err := server.ServeStdio(srv); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}