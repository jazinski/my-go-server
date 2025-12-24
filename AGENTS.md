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

## Task Tracking with Beads

**REQUIRED:** All agents MUST use Beads (`bd`) for task tracking and memory
management in every project.

### Why Beads?

- **Persistent memory**: Git-backed issue tracking that survives sessions,
  branches, and merges
- **Dependency-aware**: Agents know what tasks are ready and what's blocked
- **Zero conflicts**: Hash-based IDs (`bd-a1b2`) prevent merge collisions
- **Multi-agent safe**: Multiple agents can work in parallel without conflicts
- **Audit trail**: Every change tracked with full history

### Core Workflow

1. **Start every session** by checking ready work:
   ```bash
   bd ready --json
   ```

2. **Create beads for all work** (never use markdown TODO lists):
   ```bash
   bd create "Implement feature X" -p 1 --json
   bd create "Add tests for feature X" -p 1 --json
   ```

3. **Update status as you work**:
   ```bash
   bd update bd-abc --status in_progress --json
   ```

4. **Close completed work**:
   ```bash
   bd close bd-abc --reason "Completed" --json
   ```

5. **End of session sync**:
   ```bash
   bd sync
   ```

### Azure DevOps Integration

**IMPORTANT:** When creating beads that represent user-facing features or bugs,
also create corresponding Azure DevOps work items:

```bash
# Create bead
bd create "Fix validation bug in UserController" -t bug -p 0 --json

# Create Azure DevOps work item
az boards work-item create \
  --type Bug \
  --title "Fix validation bug in UserController" \
  --description "Linked to bead: bd-abc123" \
  --project "MyProject"

# Link them
bd update bd-abc123 --metadata "az_work_item=12345"
```

**When to create Azure DevOps work items:**

- ✅ **ALWAYS** for: User-facing features, bugs, epics, sprint commitments
- ⚠️ **OPTIONAL** for: Agent planning tasks, technical debt, internal sub-tasks
- ❌ **NEVER** for: Transient notes, duplicate tracking

### Essential Commands

```bash
# List all tasks
bd list --json

# Show task details
bd show bd-abc --json

# Add dependency (B blocks on A)
bd dep add bd-B bd-A --type blocks

# Show dependency graph
bd graph

# Sync to git
bd sync
```

### Full Documentation

See [Beads Integration Guide](assets/resources/processes/beads-integration.md)
for comprehensive workflow, patterns, and best practices.

## Team Documentation & Standards

**IMPORTANT:** This MCP server provides access to team documentation via tools.
Use these tools to load relevant standards and processes before starting work.

### Available Documentation Tools

- **`list_documentation`**: Lists all available documentation organized by
  category (coding standards, processes, architecture, examples)
- **`load_documentation`**: Loads the full content of a specific documentation
  file

### When to Load Documentation

**ALWAYS load relevant documentation before:**

- **Starting new code**: Load appropriate coding standard
  - React: `coding-standards/reactjs-style-guide.md`
  - .NET Core: `coding-standards/dotnet-core-style-guide.md`
  - ColdFusion: `coding-standards/coldfusion-style-guide.md`
  - API design: `coding-standards/api-design-guide.md`
- **Task tracking setup**: Load `processes/beads-integration.md` for Beads
  workflow
- **Azure DevOps integration**: Load `processes/azure-devops-workflow.md` for
  patterns
- **Architecture decisions**: Load `architecture/system-design-principles.md`
  for design guidance
- **Database work**: Load `coding-standards/database-conventions.md` for schema
  conventions

### Documentation Workflow Example

```bash
# 1. List available docs
list_documentation

# 2. Load relevant standard based on task
load_documentation "coding-standards/reactjs-style-guide.md"

# 3. Proceed with implementation following loaded standards
```

### Quick Reference

| Task Type              | Load This Documentation                       |
| ---------------------- | --------------------------------------------- |
| React component        | `coding-standards/reactjs-style-guide.md`     |
| .NET Core API          | `coding-standards/dotnet-core-style-guide.md` |
| Task tracking          | `processes/beads-integration.md`              |
| Azure DevOps workflows | `processes/azure-devops-workflow.md`          |
| Database schema        | `coding-standards/database-conventions.md`    |
| Git workflow           | `coding-standards/git-workflow.md`            |
| System architecture    | `architecture/system-design-principles.md`    |
| Automation examples    | `examples/beads-azure-devops-automation.md`   |

**Note:** Resources are also available via MCP resource protocol, but OpenCode
requires explicit tool calls to access them.
