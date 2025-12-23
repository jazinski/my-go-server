package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Log to stderr to avoid interfering with stdio mode JSON communication
		fmt.Fprintf(os.Stderr, "Warning: .env file not found or could not be loaded: %v\n", err)
	}

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

	// Tool to send pushover notifications
	pushTool := mcp.NewTool("send_push_notification",
		mcp.WithDescription("Send a push notification via Pushover"),
		mcp.WithString("user_key", mcp.Description("Pushover user key (optional, uses default from .env if not provided)")),
		mcp.WithString("message", mcp.Required(), mcp.Description("Notification message")),
	)
	srv.AddTool(pushTool, pushHandler)

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
		fmt.Printf("Starting MCP Server in SSE mode on :%s\n", *port)
		sseServer := server.NewSSEServer(srv)
		if err := sseServer.Start(":" + *port); err != nil {
			fmt.Printf("Error starting SSE server: %v\n", err)
		}
	} else { // Stdio mode
		// In stdio mode, we must not print anything to stdout as it breaks JSON-RPC
		// All communication must be valid JSON messages
		if err := server.ServeStdio(srv); err != nil {
			// Log errors to stderr instead of stdout
			fmt.Fprintf(os.Stderr, "Error starting stdio server: %v\n", err)
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
	// Check if Docker is available
	if _, err := exec.LookPath("docker"); err != nil {
		return mcp.NewToolResultError("Docker is not installed or not in PATH. Python execution requires Docker."), nil
	}
	
	// Make sure to get the code parameter
	code, err := request.RequireString("code")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	// Lets handle optional modules parameter
	var modules []string
	if modStr := request.GetString("modules", ""); modStr != "" {
		modules = strings.Split(modStr, ",")
	}
	// Lets create temp environment and execute the code
	tmpDir, err := os.MkdirTemp("","python_repl")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create temp directory: %v", err)), nil
	}
	defer os.RemoveAll(tmpDir) // Clean up after ourselves

	// Lets take the code and write it to a temp file
	err = os.WriteFile(path.Join(tmpDir, "script.py"), []byte(code), 0644)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write temp file: %v", err)), nil
	}

	// Lets run the code in an isolated environment
	cmdArgs := []string{
		"run",
		"--rm",
		"-v", fmt.Sprintf("%s:/app", tmpDir),
		"mcr.microsoft.com/playwright/python:v1.49.1-noble",
	}

	shArgs := []string{}

	if len(modules) > 0 {
		shArgs = append(shArgs, "python", "-m", "pip", "install", "--quiet")
		shArgs = append(shArgs, modules...)
		shArgs = append(shArgs, "&&")
	}

	shArgs = append(shArgs, "python", path.Join("app", "script.py"))	
	cmdArgs = append(cmdArgs, "sh", "-c", strings.Join(shArgs, " "))

	// Now execute the docker command with timeout
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	execCmd := exec.CommandContext(timeoutCtx, "docker", cmdArgs...)
	out, err := execCmd.CombinedOutput() // Use CombinedOutput to capture both stdout and stderr
	
	if err != nil {
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return mcp.NewToolResultError("Python code execution timed out (30s limit)"), nil
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			return mcp.NewToolResultError(fmt.Sprintf("Python code failed with exit code %d: %s", exitErr.ExitCode(), string(out))), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("Failed to execute Python code: %v\nOutput: %s", err, string(out))), nil
	}

	return mcp.NewToolResultText(string(out)), nil
}

// pushHandler is the handler for the send_push_notification tool
func pushHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get user key from request or fall back to environment variable
	userKey := request.GetString("user_key", "")
	if userKey == "" {
		userKey = os.Getenv("USER_KEY")
		if userKey == "" {
			return mcp.NewToolResultError("user_key parameter not provided and USER_KEY environment variable not set"), nil
		}
	}
	
	message, err := request.RequireString("message")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get the app token from environment variable
	appToken := os.Getenv("APP_TOKEN")
	if appToken == "" {
		return mcp.NewToolResultError("APP_TOKEN environment variable not set"), nil
	}

	// Prepare the Pushover command
	cmd := exec.Command("curl", "-s",
		"--form-string", fmt.Sprintf("token=%s", appToken),
		"--form-string", fmt.Sprintf("user=%s", userKey),
		"--form-string", fmt.Sprintf("message=%s", message),
		"https://api.pushover.net/1/messages.json",
	)

	// Execute the command
	out, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to send push notification: %v\nOutput: %s", err, string(out))), nil
	}

	return mcp.NewToolResultText("Push notification sent successfully!"), nil
}