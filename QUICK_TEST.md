# Quick Test: Beads Workflow Compliance

**5-Minute Agent Test** - Verify your AI agent follows Beads workflow

---

## 🚀 Quick Setup (30 seconds)

```bash
# Create test project
mkdir /tmp/beads-quick-test && cd /tmp/beads-quick-test
git init && bd init
```

---

## 🧪 The Test (2 minutes)

**Prompt your AI agent with:**

```
"I need to implement a user authentication feature. 
Create tasks for this work following our team's standard workflow."
```

---

## ✅ What to Look For

### Agent SHOULD do:

- ✅ Use `bd create` commands
- ✅ Create multiple beads (epic → tasks)
- ✅ NOT create markdown TODO lists

### Quick Check:

```bash
bd list  # Should show beads created
```

**PASS** if you see beads listed (bd-xxxx)\
**FAIL** if agent created markdown TODO instead

---

## 🎯 Test Azure DevOps Integration (1 minute)

**Prompt:**

```
"Create a P0 bug task for login validation issue (user-facing). 
This needs to be tracked in Azure DevOps."
```

### Agent SHOULD do:

- ✅ Create bead: `bd create "..." -t bug -p 0`
- ✅ Create Azure DevOps work item
- ✅ Link them together

### Quick Check:

```bash
bd list --json | jq '.[-1]'  # Last bead
# Check for: .metadata.az_work_item
```

**PASS** if bead has `az_work_item` in metadata\
**FAIL** if no Azure DevOps link

---

## 🎉 Results

**Both tests passed?** ✅ Your agent is following Beads workflow!

**Test failed?** ⚠️ See [TESTING_BEADS_WORKFLOW.md](TESTING_BEADS_WORKFLOW.md)
for:

- Detailed test suite
- Debugging guide
- Compliance checklist

---

## 📊 Quick Compliance Check

```bash
# Count beads created today
bd list --json | jq '. | length'

# Check for Azure DevOps links
bd list --json | jq '[.[] | select(.metadata.az_work_item != null)] | length'

# Check status distribution
bd list --json | jq 'group_by(.status) | map({status: .[0].status, count: length})'
```

---

## 🔗 Full Testing Guide

For comprehensive testing including:

- Session workflow tests
- Multi-agent coordination
- Dependency management
- Context persistence
- Automated test scripts

**See:** [TESTING_BEADS_WORKFLOW.md](TESTING_BEADS_WORKFLOW.md)

---

**Tip:** Run this quick test whenever you:

- Update the MCP server
- Configure a new AI agent
- Onboard new team members
- Change Beads configuration

---

**Last Updated:** December 23, 2025\
**Quick Test Time:** ~5 minutes\
**Comprehensive Test Time:** ~20 minutes
