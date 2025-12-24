# Team Standards MCP Server - Setup & Configuration Guide

## 🎯 Overview

This MCP server provides centralized team coding standards, workflows, and best
practices for use with OpenCode and other MCP-compatible clients.

**Key Benefits:**

- ✅ Single source of truth for team standards
- ✅ Dynamic updates (no manual sync needed)
- ✅ Always available in AI coding assistants
- ✅ Consistent development practices across team

## 📦 What's Included

### Resources (Always Available)

Located in `assets/resources/`, these are automatically loaded into the AI's
context:

- **Coding Standards:**
  - `go-style-guide.md` - Go coding conventions and best practices
  - `git-workflow.md` - Git branching, commits, PR process

- **Architecture** (planned):
  - API design guidelines
  - System design principles
  - Database conventions

- **Processes** (planned):
  - Release process
  - Incident response
  - Onboarding guide

### Prompts (On-Demand)

Located in `assets/prompts/`, invoked when needed:

- `code-review.md` - Comprehensive code review workflow
- `analyzer.md` - Content analysis template

### Tools

- `hello_tool` - Example tool (can be removed)
- `execute_python` - Python execution with Playwright support
- `send_push_notification` - Pushover notifications (optional)

## 🚀 Quick Start for Team Members

### Prerequisites

- [OpenCode](https://opencode.ai) installed
- Go 1.21+ (if running locally)
- Docker (for Python execution tool)

### Option 1: Local Server (Development)

1. **Clone the repository:**
   ```bash
   git clone <repo-url>
   cd my-go-server
   ```

2. **Build the server:**
   ```bash
   go build -o my-go-server .
   ```

3. **Add to OpenCode config:**

   Edit `~/.config/opencode/opencode.json` (or workspace
   `.opencode/opencode.json`):

   ```json
   {
     "mcpServers": {
       "team-standards": {
         "command": "/absolute/path/to/my-go-server",
         "args": [],
         "env": {}
       }
     }
   }
   ```

   Replace `/absolute/path/to/` with the actual path.

4. **Restart OpenCode** to load the server.

### Option 2: Remote Server (Production)

**For self-hosted team server:**

1. **Add to OpenCode config:**

   ```json
   {
     "mcpServers": {
       "team-standards": {
         "command": "/absolute/path/to/my-go-server",
         "args": ["-sse", "-port=8080"],
         "env": {},
         "transport": "sse",
         "url": "http://your-server:8080/sse"
       }
     }
   }
   ```

2. **Or use stdio over SSH:**

   ```json
   {
     "mcpServers": {
       "team-standards": {
         "command": "ssh",
         "args": ["user@your-server", "/path/to/my-go-server"],
         "env": {}
       }
     }
   }
   ```

## 📝 Configuration Examples

### Minimal Configuration

Bare minimum for team standards only:

```json
{
  "mcpServers": {
    "team-standards": {
      "command": "/home/user/mcp-servers/my-go-server",
      "args": []
    }
  }
}
```

### With Notifications

Enable Pushover notifications:

```json
{
  "mcpServers": {
    "team-standards": {
      "command": "/home/user/mcp-servers/my-go-server",
      "args": [],
      "env": {
        "USER_KEY": "your-pushover-user-key",
        "APP_TOKEN": "your-pushover-app-token"
      }
    }
  }
}
```

### Multiple Environments

Different configs for dev/prod:

```json
{
  "mcpServers": {
    "team-standards-dev": {
      "command": "/home/user/dev/my-go-server",
      "args": []
    },
    "team-standards-prod": {
      "command": "/home/user/prod/my-go-server",
      "args": [],
      "transport": "stdio"
    }
  }
}
```

## 🔧 Advanced Configuration

### Custom Resource Paths

If you want to organize differently:

```
assets/
├── resources/
│   ├── team-a/           # Team A specific standards
│   ├── team-b/           # Team B specific standards
│   └── shared/           # Shared standards
└── prompts/
    └── workflows/        # Workflow templates
```

The server recursively loads all `.md` files, so any structure works.

### Token Optimization

MCP resources are added to the AI's context, which uses tokens. To optimize:

1. **Keep resources concise** - Focus on essential information
2. **Use prompts for detailed workflows** - Invoked only when needed
3. **Selective loading** - Create team-specific servers with subsets of content

### Example: Team-Specific Server

```bash
# Team A uses only Go standards
ln -s assets/resources/coding-standards/go-style-guide.md assets/resources/team-a/

# Team B uses all standards
ln -s assets/resources/* assets/resources/team-b/
```

Then configure different servers in OpenCode for each team.

## 🌐 Deployment Options

### Option 1: Git Repository (Recommended)

**Pros:**

- Version controlled
- Easy updates (git pull)
- Code review for standards changes
- History tracking

**Setup:**

```bash
# Each team member:
git clone <repo-url> ~/mcp-servers/team-standards
cd ~/mcp-servers/team-standards
go build -o my-go-server .

# Add to OpenCode config (see above)

# To update:
git pull origin main
go build -o my-go-server .
# Restart OpenCode
```

### Option 2: Shared Network Drive

**Pros:**

- Automatic updates
- No git knowledge required
- Works for non-developers

**Setup:**

```bash
# Mount shared drive
sudo mount -t cifs //server/share /mnt/team-standards

# In OpenCode config:
{
  "mcpServers": {
    "team-standards": {
      "command": "/mnt/team-standards/my-go-server",
      "args": []
    }
  }
}
```

### Option 3: Self-Hosted Server

**Pros:**

- Central deployment
- No client setup beyond config
- Can add authentication

**Setup:**

1. **Server side:**
   ```bash
   cd /opt/team-standards
   ./my-go-server -sse -port=8080
   ```

2. **Optional: systemd service**
   ```ini
   [Unit]
   Description=Team Standards MCP Server
   After=network.target

   [Service]
   Type=simple
   User=mcp
   WorkingDirectory=/opt/team-standards
   ExecStart=/opt/team-standards/my-go-server -sse -port=8080
   Restart=always

   [Install]
   WantedBy=multi-user.target
   ```

3. **Client side:**
   ```json
   {
     "mcpServers": {
       "team-standards": {
         "transport": "sse",
         "url": "http://team-server:8080/sse"
       }
     }
   }
   ```

### Option 4: Docker Container

**Pros:**

- Consistent environment
- Easy deployment
- Isolated dependencies

**Dockerfile:**

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o my-go-server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/my-go-server .
COPY --from=builder /app/assets ./assets
EXPOSE 8080
CMD ["./my-go-server", "-sse", "-port=8080"]
```

**Deployment:**

```bash
docker build -t team-standards-mcp .
docker run -d -p 8080:8080 --name team-standards team-standards-mcp
```

## 🔐 Security Considerations

### Local Servers (stdio)

- ✅ No network exposure
- ✅ Runs with user's permissions
- ✅ No authentication needed
- ⚠️ Each user needs local copy

### Remote Servers (SSE)

- ⚠️ Network exposed
- ⚠️ Consider authentication (not built-in currently)
- ⚠️ Use TLS in production
- ✅ Central management

### Recommended Approach:

1. **Internal network only** - Don't expose to internet
2. **VPN required** - Team members connect via VPN
3. **Firewall rules** - Restrict to team IP ranges
4. **Future**: Add OAuth/API key authentication

## 🔄 Updating Standards

### For Repository Maintainers:

1. **Make changes to markdown files** in `assets/`
2. **Commit and push:**
   ```bash
   git add assets/
   git commit -m "docs: update Go style guide with error handling examples"
   git push origin main
   ```

3. **Notify team:**
   ```bash
   # Post in Slack/Teams:
   "📢 Team standards updated! Run `git pull` in ~/mcp-servers/team-standards 
   and restart OpenCode to get the latest."
   ```

### For Team Members:

1. **Update local copy:**
   ```bash
   cd ~/mcp-servers/team-standards
   git pull origin main
   go build -o my-go-server .
   ```

2. **Restart OpenCode** (or reload MCP servers)

### Automatic Updates (Advanced):

Create a cron job or systemd timer:

```bash
# ~/.config/systemd/user/mcp-update.service
[Unit]
Description=Update MCP Team Standards

[Service]
Type=oneshot
WorkingDirectory=%h/mcp-servers/team-standards
ExecStart=/bin/bash -c 'git pull origin main && go build -o my-go-server .'
```

```bash
# ~/.config/systemd/user/mcp-update.timer
[Unit]
Description=Update MCP Team Standards Daily

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
```

## 🧪 Testing Your Setup

After configuration, verify it works:

### 1. Check Server Starts

```bash
./my-go-server
# Should start without errors
# Press Ctrl+C to stop
```

### 2. Test in OpenCode

1. Open OpenCode
2. Check MCP servers are loaded (status bar or settings)
3. Ask AI: "What are our team's Go coding standards?"
4. AI should reference the `go-style-guide.md` resource

### 3. Test Prompts

1. Ask AI: "Run a code review on this file"
2. AI should invoke the `code-review` prompt
3. Should follow the checklist in `code-review.md`

### 4. Test Tools (Optional)

```
# Python execution:
Ask AI: "Execute Python code to print Hello World"

# Pushover notifications:
Ask AI: "Send me a push notification saying test"
```

## 🐛 Troubleshooting

### Server Won't Start

**Error: "command not found"**

```bash
# Check path in opencode.json is absolute
# Use `pwd` to get full path:
cd ~/mcp-servers/team-standards && pwd
# Use that path in opencode.json
```

**Error: "permission denied"**

```bash
chmod +x my-go-server
```

**Error: "port already in use" (SSE mode)**

```bash
# Change port in config
./my-go-server -sse -port=8081
```

### OpenCode Can't See Resources

1. **Check logs:** OpenCode usually logs MCP server output
2. **Verify files exist:**
   ```bash
   ls -R assets/
   ```
3. **Check file format:** Must be `.md` files
4. **Rebuild server:**
   ```bash
   go build -o my-go-server .
   ```

### Updates Not Showing

1. **Rebuild after changes:**
   ```bash
   go build -o my-go-server .
   ```
2. **Restart OpenCode** completely
3. **Clear cache** if OpenCode has one

### Python Tool Issues

**Error: "Docker is not installed"**

```bash
# Install Docker
sudo apt-get install docker.io   # Ubuntu/Debian
brew install docker              # macOS

# Add user to docker group
sudo usermod -aG docker $USER
# Log out and back in
```

**Playwright Issues:** See [PLAYWRIGHT.md](../../PLAYWRIGHT.md) for detailed
troubleshooting.

## 📊 Monitoring & Metrics

### Track Adoption

Create a simple analytics endpoint (optional):

```go
// Add to main.go
var usageCounter atomic.Int64

func trackUsage() {
    usageCounter.Add(1)
}

// Call trackUsage() in resource/prompt handlers
```

### Log Resource Access

Add logging to understand which standards are most referenced:

```bash
# Check OpenCode logs for MCP resource reads
tail -f ~/.config/opencode/logs/mcp.log | grep "team-standards"
```

## 🎓 Best Practices

### For Maintainers:

1. **Keep standards current** - Review quarterly
2. **Get team input** - Standards should be collaborative
3. **Use PRs for changes** - Discuss major updates
4. **Version standards** - Track changes in git history
5. **Examples over rules** - Show good code, don't just list rules
6. **Link to external docs** - Don't duplicate official documentation

### For Users:

1. **Trust the AI** - It has access to latest standards
2. **Ask questions** - "Why does our style guide recommend X?"
3. **Suggest improvements** - Create issues or PRs for outdated info
4. **Use prompts** - Leverage pre-built workflows
5. **Keep updated** - Pull latest changes regularly

## 🔮 Future Enhancements

Planned improvements:

- [ ] **Authentication** for remote servers (OAuth, API keys)
- [ ] **Custom tools** for team-specific automation
- [ ] **Multiple language standards** (Python, TypeScript, etc.)
- [ ] **Interactive tutorials** as prompts
- [ ] **AI-powered linting** custom to team standards
- [ ] **Metrics dashboard** for standards adoption
- [ ] **Auto-update mechanism** for clients

## 📚 Additional Resources

- [OpenCode Documentation](https://opencode.ai/docs)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Go MCP SDK](https://github.com/mark3labs/mcp-go)

## 🆘 Support

- **Issues:** Create an issue in the git repository
- **Questions:** Ask in team chat (Slack/Teams)
- **Urgent:** Contact DevOps team

---

**Last Updated:** 2025-01-23 **Maintainer:** [Your Name/Team] **Version:** 1.0.0
