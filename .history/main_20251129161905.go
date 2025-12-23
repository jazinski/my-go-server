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
	port := flag.String("port", "8080", "Port to run the SSE server on")
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
	// Add Tool Handler
	srv.AddTool(pyTool, pyToolHandler)

	// Start the server in the desired mode
	if *sseMode {
		fmt.Println("Starting MCP Server in SSE mode on :8080")
		sseServer := server.NewSSEServer(srv)
		if err := sseServer.Start(":8080"); err != nil {
			fmt.Printf("Error starting SSE server: %v\n", err)
		}
	} else {
		fmt.Println("Starting MCP Server in stdio mode")
		if err := server.ServeStdio(srv); err != nil {
			fmt.Printf("Error starting stdio server: %v\n", err)
		}
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

func pyToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	code, err := request.RequireString("code")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	modules, _ := request.GetString("modules")

	output, err := executePythonCode(code, modules)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(output), nil
}	