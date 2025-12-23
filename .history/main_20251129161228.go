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
	// Python executor
	pyTool := mcp.NewTool(
		"execute_python",
		mcp.WithDescription("Execute Python code in an isolated enviornment. Playwright and headless browser are available for web scraping. Use this tool when you need real time information. Only output printed to stdout or stderr is returned so ALWAYS use printstatements! Please note all code is run in an ephemeral environment and is discarded after execution so modules and code do NOT persist!"),
		mcp.WithString(
			"code", 
			mcp.Required(), 
			mcp.Description("Python code to execute")),
		mcp.WithString(
			"modules", 
			mcp.Description("Comma-separated list of Python modules your code requires. If your code required external modules you")),
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