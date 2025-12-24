# Development & Deployment Workflow

## 📋 Quick Reference

### After Every Git Pull

```bash
git pull
go build -o my-go-server .   # ⚠️ REQUIRED - binary is not in git!
# Restart OpenCode to pick up changes
```

### Making Changes

```bash
# 1. Make your changes to main.go or add files to assets/
vim main.go
vim assets/resources/coding-standards/new-standard.md

# 2. Test locally
go build -o my-go-server .
./my-go-server  # Test in stdio mode

# 3. Commit
git add .
git commit -m "feat: add new feature"

# 4. Push
git push
```

### First Time Setup (New Team Member)

```bash
# 1. Clone repository
git clone <repo-url>
cd my-go-server

# 2. Build binary (NOT in git)
go build -o my-go-server .

# 3. Configure OpenCode
# Edit ~/.config/opencode/opencode.json
{
  "mcpServers": {
    "team-standards": {
      "command": "/absolute/path/to/my-go-server/my-go-server",
      "args": []
    }
  }
}

# 4. Restart OpenCode

# 5. Test
# Ask OpenCode: "List available documentation"
```

## 🎯 Understanding What's in Git

### ✅ Checked into Git

- `main.go` - Source code
- `go.mod`, `go.sum` - Dependencies
- `assets/` - All documentation and prompts **⚠️ MUST be in same directory as
  binary**
- `README.md`, `AGENTS.md`, etc. - Documentation
- `.gitignore` - Git configuration

### ❌ NOT in Git (by design)

- `my-go-server` - Compiled binary (platform-specific)
- `python_output/` - Temporary Python execution outputs
- `.env` - Environment variables (secrets)
- `*.log` - Log files

**Why exclude the binary?**

- Platform-specific (macOS binary won't work on Linux/Windows)
- Can be large (9+ MB)
- Changes with every code update
- Easy to rebuild from source: `go build -o my-go-server .`

**⚠️ CRITICAL: Binary and Assets Must Be Together**

The binary looks for `assets/` **relative to where the binary file is located**,
not the current working directory. This means:

✅ **Correct structure:**

```
my-go-server/
├── my-go-server          (binary)
├── assets/
│   ├── resources/
│   └── prompts/
├── main.go
└── ...
```

❌ **Won't work:**

```
/usr/local/bin/my-go-server          (binary)
/somewhere/else/my-go-server/assets/  (assets)
```

**Solution:** Always keep the binary in the same directory where you cloned the
repo.

## 🔄 Common Workflows

### Scenario: Team member adds new documentation

**Person A (adds doc):**

```bash
vim assets/resources/coding-standards/python-guide.md
git add .
git commit -m "docs: add Python coding standards"
git push
```

**Person B (uses new doc):**

```bash
git pull
go build -o my-go-server .  # ⚠️ CRITICAL STEP
# Restart OpenCode
# Now can use: "Load the Python coding standards"
```

### Scenario: Adding a new tool

**Developer:**

```bash
# 1. Edit main.go to add tool
vim main.go

# 2. Test locally
go build -o my-go-server .
./my-go-server

# 3. Commit and push
git add main.go
git commit -m "feat: add database migration tool"
git push
```

**Other team members:**

```bash
git pull
go build -o my-go-server .  # ⚠️ REBUILD REQUIRED
# Restart OpenCode
# New tool now available
```

### Scenario: Multi-platform team

**macOS developer:**

```bash
git pull
go build -o my-go-server .  # Creates macOS binary
```

**Linux developer:**

```bash
git pull
go build -o my-go-server .  # Creates Linux binary
```

**Windows developer:**

```bash
git pull
go build -o my-go-server.exe  # Creates Windows binary
```

Each platform builds their own binary from the same source.

## 🐛 Common Mistakes

### ❌ Mistake #1: Forget to rebuild after pull

```bash
git pull
# Forgot: go build
# Restart OpenCode
# Result: Old version still running, new features missing!
```

**Fix:** Always rebuild after pull:

```bash
git pull
go build -o my-go-server .
```

### ❌ Mistake #2: Trying to add binary to git

```bash
git add my-go-server
# Result: Git ignores it (in .gitignore)
```

**This is correct!** Binary should NOT be in git.

### ❌ Mistake #3: Moving binary away from assets directory

```bash
# ❌ WRONG - Binary separated from assets
cp my-go-server /usr/local/bin/
# Binary can't find assets/
```

**Fix:** Keep binary in the repo directory, or copy the entire directory:

```bash
# ✅ CORRECT - Keep binary with assets
# In OpenCode config, point to binary in repo:
{
  "command": "/path/to/my-go-server/my-go-server"
}

# OR copy entire directory if needed:
cp -r my-go-server /usr/local/share/
# Then point to: /usr/local/share/my-go-server/my-go-server
```

### ❌ Mistake #4: Relative path in OpenCode config

```json
{
  "mcpServers": {
    "team-standards": {
      "command": "./my-go-server", // ❌ WRONG - relative path
      "args": []
    }
  }
}
```

**Fix:** Use absolute path:

```json
{
  "mcpServers": {
    "team-standards": {
      "command": "/Users/yourname/mcp-servers/my-go-server/my-go-server", // ✅ CORRECT
      "args": []
    }
  }
}
```

### ❌ Mistake #5: Not restarting OpenCode

```bash
git pull
go build -o my-go-server .
# Forgot to restart OpenCode
# Result: Still using old version!
```

**Fix:** Always restart OpenCode after rebuild.

## ✅ Verification Checklist

After pulling and rebuilding, verify everything works:

```bash
# 1. Check binary exists and is recent
ls -lh my-go-server
# Should show recent timestamp

# 2. Test server starts
./my-go-server
# Should not error (Ctrl+C to stop)

# 3. Restart OpenCode completely

# 4. Test in OpenCode
# Ask: "List available documentation"
# Should see all docs including any new ones

# Ask: "Load the Beads integration guide"
# Should load content successfully
```

## 🚀 Deployment Strategies

### Strategy 1: Everyone builds locally (Recommended)

- Each developer: `git pull && go build`
- Platform-independent
- Everyone on latest code

### Strategy 2: Pre-built binaries (Not recommended)

- Could store binaries in releases
- Requires separate binary per platform
- More complex to maintain

### Strategy 3: Shared network drive

- Build once on shared drive
- Team points to same binary
- Works if everyone on same OS

**Recommendation:** Strategy 1 (local builds) is simplest and most reliable.

## 📚 Related Documentation

- **[README.md](README.md)** - Main project documentation
- **[AGENTS.md](AGENTS.md)** - Agent development guidelines
- **[TROUBLESHOOTING.md](TROUBLESHOOTING.md)** - Fix common issues
- **[.gitignore](.gitignore)** - See what's excluded from git

## 🎓 Learning Resources

**New to Go?**

- Install Go: https://go.dev/doc/install
- Basic Go tutorial: https://go.dev/tour/
- All you need to know: `go build -o my-go-server .`

**New to MCP?**

- MCP Protocol: https://modelcontextprotocol.io
- OpenCode docs: https://opencode.ai/docs

---

**TL;DR:** After every `git pull`, run `go build -o my-go-server .` and restart
OpenCode!
