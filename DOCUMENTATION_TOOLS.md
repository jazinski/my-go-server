# Team Documentation Access Guide

## 🎯 Problem Solved

**Issue:** OpenCode doesn't automatically use MCP resources (passive documentation).

**Solution:** We've added **explicit tool-based documentation access** that works perfectly with OpenCode's paradigm.

---

## 🚀 Available Tools

### 1. `list_documentation`
**Purpose:** Discover all available team documentation organized by category.

**Usage:**
```bash
# Lists all documentation with categories
list_documentation
```

**Output:**
- **Coding Standards** (8 files: React, .NET, ColdFusion, API design, etc.)
- **Processes** (2 files: Beads integration, Azure DevOps workflow)
- **Architecture** (1 file: System design principles)
- **Examples** (1 file: Beads + Azure DevOps automation)

### 2. `load_documentation`
**Purpose:** Load the full content of a specific documentation file.

**Usage:**
```bash
load_documentation "processes/beads-integration.md"
```

**Parameters:**
- `path` (required): Relative path to documentation file

**Examples:**
```bash
# Load React coding standards
load_documentation "coding-standards/reactjs-style-guide.md"

# Load Beads workflow
load_documentation "processes/beads-integration.md"

# Load API design guide
load_documentation "coding-standards/api-design-guide.md"
```

---

## 📋 When to Load Documentation

### 🎨 **Starting New Code**
**Load appropriate coding standard FIRST:**

| Language/Framework | Documentation File |
|-------------------|-------------------|
| React | `coding-standards/reactjs-style-guide.md` |
| .NET Core | `coding-standards/dotnet-core-style-guide.md` |
| .NET Framework | `coding-standards/dotnet-framework-style-guide.md` |
| ColdFusion | `coding-standards/coldfusion-style-guide.md` |
| AngularJS | `coding-standards/angularjs-style-guide.md` |
| API Design | `coding-standards/api-design-guide.md` |

### 📊 **Working with Azure DevOps**
```bash
# Load ADO workflow patterns
load_documentation "processes/azure-devops-workflow.md"

# See automation examples
load_documentation "examples/beads-azure-devops-automation.md"
```

### 🎯 **Task Tracking Setup**
```bash
# Load Beads workflow guide
load_documentation "processes/beads-integration.md"
```

### 🏗️ **Architecture Decisions**
```bash
# Load system design principles
load_documentation "architecture/system-design-principles.md"
```

### 🗄️ **Database Work**
```bash
# Load database conventions
load_documentation "coding-standards/database-conventions.md"
```

### 🔀 **Git Workflow**
```bash
# Load git conventions
load_documentation "coding-standards/git-workflow.md"
```

---

## 🔄 Typical Workflow

```bash
# 1. Start session - list available docs
list_documentation

# 2. Load relevant standard for your task
load_documentation "coding-standards/reactjs-style-guide.md"

# 3. Load process guides if needed
load_documentation "processes/beads-integration.md"

# 4. Proceed with implementation following loaded standards
```

---

## 🛠️ Testing

Test the tools work correctly:

```bash
# Run test suite
python3 test_doc_tools.py
```

**Expected output:**
- ✅ `list_documentation` is registered
- ✅ `load_documentation` is registered  
- ✅ Documentation list contains all categories
- ✅ Can load specific documentation files

---

## 🎨 Architecture

### How It Works

1. **Discovery**: `list_documentation` scans `assets/resources/` directory
2. **Categorization**: Groups files by subdirectory (coding-standards, processes, etc.)
3. **Loading**: `load_documentation` reads requested markdown file
4. **Security**: Path traversal protection (`..` not allowed)

### Why This Approach?

**Previous:** Resources were loaded via MCP protocol but OpenCode ignored them.

**Now:** Tools provide **explicit, discoverable access** that OpenCode recognizes and uses.

**Benefits:**
- ✅ Works immediately with OpenCode
- ✅ Auto-discovers new documentation
- ✅ Clean interface (2 tools instead of 13+)
- ✅ Searchable and browsable
- ✅ Security-conscious (path validation)
- ✅ Maintains existing resource system for future compatibility

---

## 📝 Adding New Documentation

Simply add markdown files to `assets/resources/`:

```bash
assets/resources/
├── coding-standards/
│   └── your-new-guide.md          # Auto-discovered
├── processes/
│   └── your-new-process.md        # Auto-discovered
├── architecture/
│   └── your-new-pattern.md        # Auto-discovered
└── examples/
    └── your-new-example.md        # Auto-discovered
```

**No code changes needed!** Tools auto-discover new files.

---

## 🔍 Available Documentation

### Coding Standards (8 files)
- `angularjs-style-guide.md` - AngularJS coding conventions
- `api-design-guide.md` - RESTful API design patterns
- `coldfusion-style-guide.md` - ColdFusion best practices
- `database-conventions.md` - Database schema and naming
- `dotnet-core-style-guide.md` - .NET Core standards
- `dotnet-framework-style-guide.md` - .NET Framework patterns
- `git-workflow.md` - Git branching and commit conventions
- `reactjs-style-guide.md` - React component standards

### Processes (2 files)
- `azure-devops-workflow.md` - ADO integration patterns
- `beads-integration.md` - Task tracking with Beads

### Architecture (1 file)
- `system-design-principles.md` - System architecture guide

### Examples (1 file)
- `beads-azure-devops-automation.md` - Automation patterns

---

## 🚨 Troubleshooting

### "Failed to load documentation"
**Cause:** File path incorrect or file doesn't exist

**Fix:** Use `list_documentation` first to see exact file paths

### "Invalid path: path traversal not allowed"
**Cause:** Path contains `..` (security protection)

**Fix:** Use relative paths from `assets/resources/` only

### Tools not showing up
**Cause:** MCP server not rebuilt after changes

**Fix:**
```bash
go build -o my-go-server .
# Restart your MCP client (OpenCode, etc.)
```

---

## 🎯 Key Differences from Resources

| Feature | MCP Resources (Old) | Documentation Tools (New) |
|---------|-------------------|-------------------------|
| **Discoverability** | Hidden from OpenCode | Explicit tool list |
| **Access** | Passive (unused) | Active (function calls) |
| **Auto-discovery** | ❌ | ✅ |
| **OpenCode Support** | ❌ | ✅ |
| **Security** | N/A | Path validation |
| **Categorization** | Manual | Automatic |

---

## 📚 Related Documentation

- **AGENTS.md** - Agent guidelines and team standards section
- **README.md** - Main MCP server documentation  
- **QUICK_TEST.md** - Testing documentation

---

## ✅ Success Criteria

You'll know it's working when:

1. ✅ `list_documentation` shows all your team docs categorized
2. ✅ `load_documentation` returns full markdown content
3. ✅ Agents reference loaded docs in their work
4. ✅ New docs appear automatically without code changes
5. ✅ OpenCode recognizes and uses the tools

---

**Status:** ✅ **IMPLEMENTED & TESTED**\
**Version:** 1.0\
**Last Updated:** December 2025
