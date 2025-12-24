# Team Standards MCP Server

> **Centralized coding standards, workflows, and best practices for AI-powered
> development teams**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![MCP Protocol](https://img.shields.io/badge/MCP-Compatible-green)](https://modelcontextprotocol.io)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 🎯 What is This?

This is an **MCP (Model Context Protocol) server** that provides your
development team with:

- 📚 **Centralized Coding Standards**: Go style guides, Git workflows, and more
- 🤖 **AI-Ready**: Works with OpenCode, Claude Code, and other MCP-compatible AI
  assistants
- 🔄 **Always Up-to-Date**: Update standards in one place, all team members get
  changes
- 🎓 **Built-in Workflows**: Pre-built prompts for code reviews, releases, and
  common tasks
- 🛠️ **Extensible**: Add custom tools, resources, and prompts for your team

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** (for building)
- **OpenCode** or another MCP-compatible client
- **Docker** (optional, for Python execution tool)

### Installation

1. **Clone and build:**
   ```bash
   git clone <your-repo-url>
   cd my-go-server
   go build -o my-go-server .
   ```

2. **Configure OpenCode:**

   Add to `~/.config/opencode/opencode.json`:
   ```json
   {
     "mcpServers": {
       "team-standards": {
         "command": "/absolute/path/to/my-go-server",
         "args": []
       }
     }
   }
   ```

3. **Restart OpenCode** and you're done! 🎉

### Verify It Works

Ask your AI assistant:

```
"What are our team's Go coding standards?"
"Run a code review on this code"
"Explain our Git workflow"
```

The AI will reference your team standards automatically!

## 📚 What's Included

### Resources (Always Available in AI Context)

#### Coding Standards

- **[.NET Core Style Guide](assets/resources/coding-standards/dotnet-core-style-guide.md)** -
  Modern C# conventions
  - Naming, formatting, async/await patterns
  - Entity Framework Core, dependency injection
  - xUnit testing and best practices
  - Security guidelines and OWASP compliance

- **[ReactJS Style Guide](assets/resources/coding-standards/reactjs-style-guide.md)** -
  Modern React patterns
  - Hooks, functional components, TypeScript
  - State management and Context API
  - Testing with React Testing Library
  - Accessibility and performance optimization

- **[.NET Framework Style Guide](assets/resources/coding-standards/dotnet-framework-style-guide.md)** -
  Legacy .NET maintenance
  - ASP.NET MVC 5 and Web Forms patterns
  - Entity Framework 6, Web API 2
  - Security and performance best practices
  - Migration planning to .NET Core

- **[AngularJS Style Guide](assets/resources/coding-standards/angularjs-style-guide.md)** -
  Legacy AngularJS 1.x maintenance
  - Component patterns and best practices
  - Controllers, services, and directives
  - Testing and performance optimization
  - Migration considerations

- **[ColdFusion Style Guide](assets/resources/coding-standards/coldfusion-style-guide.md)** -
  CFML maintenance
  - CFScript vs tag-based syntax
  - Component-based architecture (CFCs)
  - Security and parameterized queries
  - Legacy code maintenance

- **[API Design Guidelines](assets/resources/coding-standards/api-design-guide.md)** -
  REST API standards
  - Resource-oriented design and URL conventions
  - HTTP methods, status codes, and error handling
  - API versioning strategies (URL, header, query)
  - Authentication, authorization, and rate limiting
  - Pagination, filtering, and sorting
  - OpenAPI/Swagger documentation

- **[Database Conventions](assets/resources/coding-standards/database-conventions.md)** -
  Database design standards
  - Naming conventions (PascalCase, singular table names)
  - Entity Framework Core best practices
  - Migrations and schema management
  - Indexes, relationships, and query performance
  - Security (parameterized queries, connection strings)
  - Legacy database considerations

- **[Git Workflow](assets/resources/coding-standards/git-workflow.md)** - Git
  with Azure DevOps
  - Branch naming with AB# work item linking
  - Commit message conventions (Conventional Commits)
  - Pull request process
  - Code review guidelines

#### Processes & Workflows

- **[Azure DevOps Workflow](assets/resources/processes/azure-devops-workflow.md)** -
  Team processes
  - Work item types (Epic, Feature, User Story, Task, Bug)
  - Sprint planning and daily standups
  - PR workflow and branch strategy
  - Release process
  - 🤖 AI agent integration with Beads

- **[Beads Integration Guide](assets/resources/processes/beads-integration.md)** -
  **REQUIRED** task tracking for AI agents
  - Git-backed issue tracker with persistent memory
  - Dependency-aware task graphs
  - Zero-conflict merges for multi-agent workflows
  - Azure DevOps integration patterns
  - Complete agent workflow and best practices

#### Examples & Automation

- **[Beads + Azure DevOps Automation](assets/resources/examples/beads-azure-devops-automation.md)** -
  Practical automation scripts
  - Automatic work item creation hooks
  - Bidirectional status synchronization
  - Agent session management
  - Multi-agent coordination examples
  - Bash and PowerShell scripts

#### Architecture

- **[System Design Principles](assets/resources/architecture/system-design-principles.md)** -
  Core architectural principles
  - Simplicity First, Separation of Concerns, Dependency Injection
  - Fail Fast/Fail Loud, Explicit Over Implicit, Design for Testability
  - Configuration Over Code, Graceful Degradation, Security by Design
  - Monitoring & Observability
  - Comprehensive .NET Core and React examples
  - Anti-patterns to avoid
  - Architecture Decision Records (ADR) template

### Prompts (Invoked On-Demand)

- **[Code Review](assets/prompts/code-review.md)** - Comprehensive review
  checklist
  - Code quality and style checks
  - Security review (OWASP Top 10)
  - Testing validation
  - Performance analysis

- **[Accessibility Expert](assets/prompts/accessibility-expert.md)** - WCAG 2.2
  compliance
  - Semantic HTML and ARIA patterns
  - Keyboard navigation and screen readers
  - Color contrast and form accessibility
  - Testing checklist

- **[Security Expert](assets/prompts/security-expert.md)** - Application
  security
  - OWASP Top 10 (2021) prevention
  - SQL injection, XSS, CSRF protection
  - Authentication and authorization patterns
  - Secure coding for C# and React

- **[UI/UX Expert](assets/prompts/analyzer.md)** - Design and usability
  - Visual hierarchy and consistency
  - Feedback and affordance patterns
  - Responsive design principles
  - UX review checklist

### Tools

- **`execute_python`** - Run Python code in isolated Docker container
  - Playwright support for web scraping
  - 30-second timeout
  - Output files persist to `python_output/`

- **`send_push_notification`** - Send Pushover notifications (optional)

- **`hello_tool`** - Example tool (can be removed)

## 🏗️ Architecture

```
my-go-server/
├── main.go                    # MCP server implementation
├── assets/
│   ├── resources/             # Auto-loaded into AI context
│   │   ├── coding-standards/
│   │   │   ├── dotnet-core-style-guide.md
│   │   │   ├── reactjs-style-guide.md
│   │   │   ├── api-design-guide.md
│   │   │   ├── database-conventions.md
│   │   │   ├── dotnet-framework-style-guide.md
│   │   │   ├── angularjs-style-guide.md
│   │   │   ├── coldfusion-style-guide.md
│   │   │   └── git-workflow.md
│   │   ├── architecture/
│   │   │   └── system-design-principles.md
│   │   ├── processes/
│   │   │   ├── azure-devops-workflow.md
│   │   │   └── beads-integration.md         # NEW: Required for agents
│   │   └── examples/
│   │       └── beads-azure-devops-automation.md  # NEW: Automation scripts
│   └── prompts/               # Invoked on-demand
│       ├── code-review.md
│       ├── accessibility-expert.md
│       ├── security-expert.md
│       └── analyzer.md
├── python_output/             # Output from Python execution
├── go.mod
└── go.sum
```

### How It Works

1. **Server starts** and recursively loads all `.md` files from `assets/`
2. **Resources** are registered with MCP and made available to AI
3. **Prompts** are registered for on-demand invocation
4. **Tools** are registered for code execution capabilities
5. **OpenCode** (or other MCP client) connects and can access all of the above

## 🔧 Configuration

### Basic Setup (stdio mode)

```json
{
  "mcpServers": {
    "team-standards": {
      "command": "/path/to/my-go-server",
      "args": []
    }
  }
}
```

### Remote Server (SSE mode)

```json
{
  "mcpServers": {
    "team-standards": {
      "command": "/path/to/my-go-server",
      "args": ["-sse", "-port=8080"],
      "transport": "sse",
      "url": "http://localhost:8080/sse"
    }
  }
}
```

### With Environment Variables

```json
{
  "mcpServers": {
    "team-standards": {
      "command": "/path/to/my-go-server",
      "args": [],
      "env": {
        "USER_KEY": "your-pushover-user-key",
        "APP_TOKEN": "your-pushover-app-token"
      }
    }
  }
}
```

For complete configuration options, see
[TEAM_STANDARDS_SETUP.md](TEAM_STANDARDS_SETUP.md).

## 📖 Usage Examples

### Example 1: Code Review

**Developer:**

```
"Review this React component for accessibility and our team standards"
```

**AI:** _(Invokes accessibility-expert.md and checks against
reactjs-style-guide.md)_

```
I'll conduct a comprehensive accessibility review following your team standards...

🚨 MUST FIX:
- Line 23: Button is missing accessible label. Add aria-label or visible text.
- Line 45: Form input lacks associated label element (required per WCAG 2.2 AA)

⚠️ SHOULD FIX:
- Consider using semantic <button> instead of <div onClick>
- Component name should be PascalCase per your React style guide

✅ PRAISE:
- Excellent use of React.memo for performance optimization!
- Proper keyboard navigation with onKeyDown handlers
```

### Example 2: Azure DevOps Workflow Question

**Developer:**

```
"What's our branch naming convention for a new user authentication feature 
linked to work item 1234?"
```

**AI:** _(References azure-devops-workflow.md and git-workflow.md)_

```
Based on your team's Azure DevOps workflow, you should create a feature branch:

Branch name: feature/AB#1234-user-authentication

Commands:
git checkout main
git pull origin main
git checkout -b feature/AB#1234-user-authentication

This follows your team's convention:
- Type: feature/ (for new features)
- Work item: AB#1234 (auto-links to Azure DevOps)
- Description: user-authentication (lowercase, hyphens, descriptive)

When you commit, use format:
feat(auth): add OAuth2 provider [AB#1234]
```

### Example 3: Writing New Code

**Developer:**

```
"Write a C# service method to validate email addresses following our team standards"
```

**AI:** _(Uses dotnet-core-style-guide.md conventions)_

```csharp
public class EmailValidationService : IEmailValidationService
{
    private readonly ILogger<EmailValidationService> _logger;
    private static readonly Regex EmailRegex = new Regex(
        @"^[^@\s]+@[^@\s]+\.[^@\s]+$",
        RegexOptions.Compiled | RegexOptions.IgnoreCase
    );

    public EmailValidationService(ILogger<EmailValidationService> logger)
    {
        _logger = logger ?? throw new ArgumentNullException(nameof(logger));
    }

    /// <summary>
    /// Validates an email address format.
    /// </summary>
    /// <param name="email">The email address to validate.</param>
    /// <returns>True if valid, false otherwise.</returns>
    public bool IsValidEmail(string? email)
    {
        if (string.IsNullOrWhiteSpace(email))
        {
            _logger.LogWarning("Email validation failed: null or empty");
            return false;
        }

        var isValid = EmailRegex.IsMatch(email);
        if (!isValid)
        {
            _logger.LogDebug("Email validation failed for: {Email}", email);
        }

        return isValid;
    }
}
```

_This follows your team's .NET Core standards:_

- _PascalCase for public methods_
- _Nullable reference types enabled_
- _Dependency injection for ILogger_
- _XML documentation comments_
- _Structured logging with parameters_

## 🎨 Customization

### Adding New Standards

1. **Create a markdown file:**
   ```bash
   # Create new standard document
   vim assets/resources/coding-standards/api-design-guide.md
   ```

2. **Write your standards:**
   ```markdown
   # API Design Standards

   ## Overview

   Our team's REST API design conventions...

   ## Naming Conventions

   - Use plural nouns for resources (/customers, not /customer)
   - Use kebab-case for multi-word resources...
   ```

3. **Rebuild and restart:**
   ```bash
   go build -o my-go-server .
   # Restart OpenCode
   ```

4. **Done!** AI now knows your Python standards.

### Adding New Prompts

1. **Create prompt file:**
   ```bash
   vim assets/prompts/release-checklist.md
   ```

2. **Write workflow:**
   ```markdown
   # Release Preparation Workflow

   Follow these steps to prepare a release:

   1. Update version numbers...
   2. Update CHANGELOG...
   3. Run full test suite... ...
   ```

3. **Rebuild** and the prompt is available.

### Adding Custom Tools

Edit `main.go` and add your tool:

```go
// Add your custom tool
customTool := mcp.NewTool("custom_tool",
    mcp.WithDescription("Your custom tool description"),
    mcp.WithString("param", mcp.Required(), mcp.Description("Parameter description")),
)

srv.AddTool(customTool, customToolHandler)

func customToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    param, err := request.RequireString("param")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }
    
    // Your tool logic here
    result := doSomething(param)
    
    return mcp.NewToolResultText(result), nil
}
```

## 🚀 Deployment Options

### Option 1: Git Repository (Recommended)

- Each team member clones and builds locally
- `git pull` to update standards
- Version controlled, code-reviewed changes

### Option 2: Shared Network Drive

- Build once, share binary
- Automatic updates for all team members
- No git knowledge required

### Option 3: Self-Hosted Server

- Run in SSE mode on central server
- Team members connect via URL
- Single deployment, central control

### Option 4: Docker Container

- Consistent environment
- Easy deployment with orchestration
- Isolated dependencies

See [TEAM_STANDARDS_SETUP.md](TEAM_STANDARDS_SETUP.md) for detailed deployment
instructions.

## 🔐 Security

### Local (stdio mode)

- ✅ No network exposure
- ✅ Runs with user permissions
- ✅ Private to each developer

### Remote (SSE mode)

- ⚠️ Network exposed - use internal network only
- ⚠️ No built-in authentication (add reverse proxy if needed)
- ⚠️ Use TLS in production

**Recommendation:** Use stdio mode for individual developers, SSE for controlled
internal networks.

## 🧪 Testing

### Run Tests

```bash
go test ./...
```

_(Currently no tests - contributions welcome!)_

### Manual Testing

```bash
# Test server starts
./my-go-server

# Test with OpenCode
# Ask AI to reference team standards

# Test Python tool
# Ask AI to execute Python code

# Test prompts
# Ask AI to run a code review
```

## 🛠️ Development

### Build Commands

```bash
# Build
go build -o my-go-server .

# Run (stdio mode)
./my-go-server

# Run (SSE mode)
./my-go-server -sse -port=8080

# Format code
go fmt ./...

# Lint
go vet ./...
```

### Code Style

This project's server implementation is written in Go and follows standard Go
conventions.

Key points for server development:

- Use `gofmt` for formatting
- Handle all errors
- Use `context.Context` for I/O operations
- Document exported functions
- See [AGENTS.md](AGENTS.md) for detailed guidelines

**Note:** The coding standards served by this MCP server cover .NET Core, React,
AngularJS, .NET Framework, and ColdFusion - not Go. The server itself is just
the delivery mechanism.

## 🤝 Contributing

We welcome contributions! Here's how:

1. **Fork the repository**
2. **Create a feature branch:**
   ```bash
   git checkout -b feature/add-python-standards
   ```
3. **Add your standards/prompts/tools**
4. **Test with OpenCode**
5. **Commit using Conventional Commits:**
   ```bash
   git commit -m "docs: add Python coding standards"
   ```
6. **Create a Pull Request**

See our [Git Workflow Guide](assets/resources/coding-standards/git-workflow.md)
for details.

## 📊 Token Usage Considerations

MCP resources are added to the AI's context, consuming tokens. Tips for
optimization:

1. **Keep resources concise** - Focus on essential guidelines
2. **Use prompts for details** - Load only when needed
3. **Organize by team** - Create team-specific subsets
4. **Monitor context size** - OpenCode shows token usage

Example token usage:

- .NET Core Style Guide: ~3,000 tokens
- ReactJS Style Guide: ~3,000 tokens
- API Design Guidelines: ~7,000 tokens
- Database Conventions: ~6,000 tokens
- System Design Principles: ~11,000 tokens
- Azure DevOps Workflow: ~2,500 tokens
- Beads Integration Guide: ~4,000 tokens
- Beads Automation Examples: ~2,800 tokens
- Git Workflow: ~2,000 tokens
- Legacy style guides: ~2,500 tokens each (AngularJS, .NET Framework,
  ColdFusion)
- Security Expert Prompt: ~2,500 tokens (only when invoked)
- Accessibility Expert Prompt: ~2,000 tokens (only when invoked)

**Total context overhead:** ~38,000-46,000 tokens (resources only, depending on
active tech stack)

**Note:** Modern development (Core/.NET + React + API + DB + Architecture +
Beads) uses ~38,000 tokens. Legacy guides are only loaded when relevant.

## 🔮 Roadmap

- [x] .NET Core coding standards
- [x] ReactJS coding standards
- [x] Azure DevOps workflow documentation
- [x] Legacy tech standards (AngularJS, .NET Framework, ColdFusion)
- [x] Accessibility expert prompt (WCAG 2.2)
- [x] Security expert prompt (OWASP Top 10)
- [x] UI/UX expert prompt
- [x] API design guidelines
- [x] Database conventions and patterns
- [x] System design principles (comprehensive .NET/React examples)
- [x] **Beads integration for AI agent memory** (Phase 6) ✨
  - Git-backed task tracking with persistent memory
  - Azure DevOps integration patterns
  - Automation scripts and multi-agent coordination
- [ ] Authentication for remote servers
- [ ] Metrics/analytics dashboard
- [ ] Auto-update mechanism
- [ ] Custom linting tools integration
- [ ] Team-specific configuration profiles

## 📚 Resources

### Documentation

- **[Setup Guide](TEAM_STANDARDS_SETUP.md)** - Detailed configuration and
  deployment
- **[Agent Guidelines](AGENTS.md)** - Development guidelines for this project
- **[Playwright Guide](PLAYWRIGHT.md)** - Python tool with Playwright
- **[Quick Test](QUICK_TEST.md)** - 5-minute Beads workflow test ✨
- **[Testing Guide](TESTING_BEADS_WORKFLOW.md)** - Comprehensive test suite ✨
- **[Troubleshooting](TROUBLESHOOTING.md)** - Fix common Beads issues ✨

### External Resources

- [OpenCode Documentation](https://opencode.ai/docs)
- [Model Context Protocol](https://modelcontextprotocol.io)
- [MCP Go SDK](https://github.com/mark3labs/mcp-go)
- [Conventional Commits](https://www.conventionalcommits.org/)

## 🐛 Troubleshooting

### Server won't start

```bash
# Check it builds
go build -o my-go-server .

# Check file permissions
chmod +x my-go-server

# Check path in opencode.json is absolute
```

### Resources not loading

```bash
# Verify files exist
ls -R assets/

# Check they're .md files
# Rebuild after adding files
go build -o my-go-server .
```

### OpenCode can't see standards

1. Check OpenCode MCP server status
2. Verify path in config is correct
3. Restart OpenCode completely
4. Check OpenCode logs for errors

For more troubleshooting, see
[TEAM_STANDARDS_SETUP.md](TEAM_STANDARDS_SETUP.md).

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [mcp-go](https://github.com/mark3labs/mcp-go) by Mark3 Labs
- Inspired by [OpenCode](https://opencode.ai) and the MCP ecosystem
- Thanks to all contributors and the AI coding tools community

## 📧 Support

- **Issues:** [Create an issue](https://github.com/your-org/your-repo/issues)
- **Discussions:**
  [GitHub Discussions](https://github.com/your-org/your-repo/discussions)
- **Email:** team@yourcompany.com

---

**Made with ❤️ for .NET and React development teams using AI coding assistants**

_Last updated: 2025-12-23_
