package main

import (
	"context"
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

// helloHandler is the handler for the hello_tool
func helloHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := req.Params["name"].(string)
	return &mcp.ToolResponse{
		Result: fmt.Sprintf("Hello, %s!", name),
	}, nil
}