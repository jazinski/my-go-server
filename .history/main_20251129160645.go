package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Do we want SSE? (HTTP Server) or stdio?
	sseMode := flag.Bool("sse", false, "Run in SSE mode instead of stdio")
	flag.Parse()

	// Create MCP Server
	srv := server.NewMCPServer(
		"HelloMCPServer",
		"1.0.0",
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
func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}