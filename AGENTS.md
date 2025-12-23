# Agent Guidelines for my-go-server

## Build & Run Commands

- **Build**: `go build -o my-go-server .`
- **Run (stdio)**: `./my-go-server` or `go run main.go`
- **Run (SSE mode)**: `./my-go-server -sse -port=8080`
- **Test**: `go test ./...` (no tests currently exist)
- **Format**: `go fmt ./...`
- **Lint**: `go vet ./...`

## Code Style

- **Imports**: Standard library first, then third-party (godotenv, mcp-go),
  grouped with blank lines
- **Formatting**: Use `gofmt` standard formatting, tabs for indentation
- **Naming**: camelCase for private, PascalCase for exported; descriptive names
  (e.g., `helloHandler`, `pyToolHandler`)
- **Error Handling**: Always check errors; return `mcp.NewToolResultError()` for
  tool handlers, log to stderr in main
- **Logging**: Use `fmt.Fprintf(os.Stderr, ...)` to avoid breaking stdio mode
  JSON-RPC communication
- **Context**: Accept `context.Context` as first parameter in handlers
- **Tool Handlers**: Signature
  `func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)`
- **Required Params**: Use `request.RequireString()` for required;
  `request.GetString(key, default)` for optional
- **Cleanup**: Use `defer` for cleanup (e.g., `defer os.RemoveAll(tmpDir)`)

## Architecture

- MCP server with tools, resources (markdown from `assets/resources/`), and
  prompts (markdown from `assets/prompts/`)
- Supports both stdio (default) and SSE modes via `-sse` flag
- Tools use Docker for isolated Python execution; require Docker installed

## Python Execution & Playwright

- Uses Docker image: `mcr.microsoft.com/playwright/python:v1.57.0-noble`
- **Playwright version is automatically pinned to 1.57.0** when requested as a
  module
- Simply request `modules="playwright"` - no version specification needed
- Screenshots and files must be saved to `/output` (maps to `python_output/` on
  host)
- See [PLAYWRIGHT.md](PLAYWRIGHT.md) for detailed usage examples and
  troubleshooting
