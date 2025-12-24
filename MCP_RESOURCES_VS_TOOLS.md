# MCP Resources vs Tools: OpenCode Compatibility Guide

## 🎯 The Core Issue

**Question:** "OpenCode doesn't use the resources found in this MCP server. Is this a funky OpenCode-only thing?"

**Answer:** Yes and no. OpenCode follows the MCP spec correctly, but resources are **passive** and require client support to be useful.

---

## 📘 Understanding MCP Components

### Resources (Passive Documentation)
- **What:** Content that clients *can* read if they want to
- **Access:** Via `resources/list` and `resources/read` MCP protocol
- **Behavior:** "Here's documentation if you need it"
- **OpenCode Support:** ❌ Not automatically loaded or used

### Tools (Active Functions)  
- **What:** Callable functions that perform actions
- **Access:** Via `tools/list` and `tools/call` MCP protocol
- **Behavior:** "Execute this function and get results"
- **OpenCode Support:** ✅ Fully supported and actively used

### Prompts (Templates)
- **What:** Pre-configured prompt templates
- **Access:** Via `prompts/list` and `prompts/get` MCP protocol
- **Behavior:** "Use this prompt template"
- **OpenCode Support:** ⚠️ Partially (manual invocation needed)

---

## 🔍 Your Options Comparison

### Option 1: Create Explicit Documentation Tools ⭐⭐⭐⭐⭐
**Status:** ✅ IMPLEMENTED

```go
// Two simple tools for discovery and loading
addSmartResourceTools(srv)
```

**Pros:**
- ✅ Works immediately with OpenCode
- ✅ Auto-discovers new docs
- ✅ Clean interface (2 tools)
- ✅ Searchable and browsable
- ✅ Security-conscious

**Cons:**
- ⚠️ Requires two tool calls (list, then load)

**Best for:** Production use, OpenCode compatibility, maintainability

---

### Option 2: Individual Tool per Document ⭐⭐⭐
**Status:** ❌ Not implemented (available if needed)

```go
// Create separate tool for each doc
"load_reactjs_guide"
"load_beads_integration"
"load_azure_devops_workflow"
// ... (13+ tools)
```

**Pros:**
- ✅ One-call access
- ✅ Very explicit naming

**Cons:**
- ❌ Clutters tool list (13+ tools)
- ❌ Manual updates when adding docs
- ❌ Not scalable

**Best for:** Small, stable documentation sets

---

### Option 3: Embed in CLI Config Files ⭐⭐
**Status:** ❌ Not implemented (OpenCode-specific)

```bash
# Copy to OpenCode's config directory
cp -r assets/resources/ ~/.config/opencode/cli/team-docs/
```

**Pros:**
- ✅ No MCP server changes
- ✅ Always loaded in context

**Cons:**
- ❌ OpenCode-specific (not portable)
- ❌ May not work across versions
- ❌ Less discoverable
- ❌ Bloats context unnecessarily

**Best for:** Personal setups, quick testing

---

### Option 4: Keep Resources Only ⭐
**Status:** ✅ Already exists (but unused by OpenCode)

```go
// Your current implementation
loadResourcesAndPrompts(srv)
```

**Pros:**
- ✅ Follows MCP spec correctly
- ✅ Works with resource-aware clients

**Cons:**
- ❌ OpenCode ignores them
- ❌ Not discoverable
- ❌ Passive (unused)

**Best for:** Future compatibility, other MCP clients

---

### Option 5: Hybrid Approach (RECOMMENDED) ⭐⭐⭐⭐⭐
**Status:** ✅ IMPLEMENTED

```go
// Keep resources for future compatibility
loadResourcesAndPrompts(srv)

// Add tools for immediate OpenCode use
addSmartResourceTools(srv)

// Guide agents via AGENTS.md
// (documentation workflow instructions)
```

**Pros:**
- ✅ Works now with OpenCode (tools)
- ✅ Future-proof (resources still there)
- ✅ Best of both worlds
- ✅ Automatic discovery
- ✅ Agent guidance included

**Cons:**
- None significant

**Best for:** Production use (this is what we implemented)

---

## 🎨 What We Implemented

### Smart Discovery Tools

**Two tools that solve the problem elegantly:**

1. **`list_documentation`**
   - Scans `assets/resources/` directory
   - Groups by category automatically
   - Returns formatted list

2. **`load_documentation`**  
   - Loads specific file by path
   - Security: Prevents path traversal
   - Returns full markdown content

### Updated Agent Guidance

**AGENTS.md now includes:**
- When to load docs (with examples)
- Quick reference table
- Workflow patterns
- Tool usage instructions

### Testing

**Comprehensive test suite:**
```bash
python3 test_doc_tools.py
```

All tests passing ✅

---

## 🚀 Usage Examples

### Agent Workflow

```bash
# 1. Discover available docs
list_documentation

# Output:
# 📚 Available Documentation:
# **Coding Standards:**
#   - coding-standards/reactjs-style-guide.md
#   - coding-standards/dotnet-core-style-guide.md
#   ...

# 2. Load relevant standard
load_documentation "coding-standards/reactjs-style-guide.md"

# 3. Proceed with implementation following loaded standards
```

### Adding New Documentation

```bash
# Just add markdown files - auto-discovered!
echo "# New Guide" > assets/resources/processes/new-process.md

# Rebuild server
go build -o my-go-server .

# Done! Now appears in list_documentation
```

---

## 🎯 Why OpenCode Doesn't Use Resources

### Technical Reasons

1. **Protocol Design:** Resources are optional in MCP spec
2. **Client Choice:** Clients decide what to use
3. **Context Management:** Loading all resources bloats context
4. **Explicit vs Implicit:** Tools are explicit actions, resources are implicit data

### Practical Reasons

1. **Discoverability:** Agents need to know resources exist
2. **Control:** Tools give agents explicit control over loading
3. **Efficiency:** Load only what's needed when needed
4. **Workflow:** Tools fit natural agent workflow patterns

---

## 📊 Comparison Matrix

| Feature | Resources | Individual Tools | Smart Tools (✅ Implemented) |
|---------|-----------|------------------|---------------------------|
| **OpenCode Support** | ❌ | ✅ | ✅ |
| **Auto-Discovery** | ❌ | ❌ | ✅ |
| **Scalability** | N/A | ❌ | ✅ |
| **Tool List Size** | N/A | 13+ tools | 2 tools |
| **Maintenance** | Easy | Hard | Easy |
| **Security** | N/A | ✅ | ✅ |
| **Categorization** | Manual | Manual | Automatic |
| **Future-Proof** | ✅ | ❌ | ✅ |

---

## 🔮 Future Considerations

### If OpenCode Adds Resource Support

Your resources are still there! No changes needed - they'll just start working automatically.

### If You Switch MCP Clients

- **Claude Desktop:** Uses resources ✅
- **Continue.dev:** Uses resources ✅  
- **Custom clients:** Depend on implementation

Our hybrid approach works with all of them.

---

## 📝 Summary

### What Changed

**Before:**
- ✅ Resources loaded (lines 308-393 in main.go)
- ❌ OpenCode ignored them
- ❌ No way to access team docs

**After:**
- ✅ Resources still loaded (future compatibility)
- ✅ Smart tools added (immediate OpenCode use)
- ✅ Agent guidance in AGENTS.md
- ✅ Comprehensive testing
- ✅ Full documentation

### What You Get

1. **Immediate Solution:** Tools work with OpenCode now
2. **Future-Proof:** Resources ready if OpenCode adds support
3. **Maintainable:** Auto-discovery means no manual updates
4. **Scalable:** Add docs without code changes
5. **Tested:** Comprehensive test suite included
6. **Documented:** Full usage guide for agents

---

## ✅ Success Metrics

- [x] Tools registered and discoverable
- [x] OpenCode can list documentation
- [x] OpenCode can load documentation
- [x] Auto-discovers new files
- [x] Security validated (path traversal blocked)
- [x] Tests passing (3/3)
- [x] Documentation complete
- [x] Agent guidance added to AGENTS.md

---

## 📚 Related Files

- **main.go** - Implementation (lines 1-90: new tools added)
- **AGENTS.md** - Agent usage guide (updated)
- **DOCUMENTATION_TOOLS.md** - Detailed tool documentation
- **test_doc_tools.py** - Comprehensive test suite
- **assets/resources/** - Your actual documentation files

---

**Status:** ✅ **FULLY IMPLEMENTED & TESTED**\
**Version:** 1.0\
**Solution:** Hybrid approach (Resources + Smart Tools)\
**Last Updated:** December 2025
