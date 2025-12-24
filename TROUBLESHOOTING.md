# Beads Workflow Troubleshooting Guide

**Quick reference for fixing common Beads workflow issues**

---

## 🔍 Problem: Agent Not Using Beads

### Symptom

Agent creates markdown TODO lists instead of using `bd` commands

### Diagnosis

```bash
# Check if MCP server is running
# In OpenCode: Check MCP servers panel

# Check if AGENTS.md is loaded
# Ask agent: "What is our required task tracking tool?"
# Expected answer: "Beads (bd)"
```

### Solutions

**1. MCP Server Not Connected**

```bash
# Restart OpenCode
# Check ~/.config/opencode/opencode.json

# Verify server path is correct
{
  "mcpServers": {
    "team-standards": {
      "command": "/home/cjazinski/mcp-servers/my-go-server/my-go-server",
      "args": []
    }
  }
}
```

**2. Agent Not Reading AGENTS.md**

```bash
# Verify file exists
ls -l /home/cjazinski/mcp-servers/my-go-server/AGENTS.md

# Rebuild server
cd /home/cjazinski/mcp-servers/my-go-server
go build -o my-go-server .

# Restart OpenCode
```

**3. Agent Needs Explicit Instruction**

```
# Add to prompt:
"Use Beads (bd) for task tracking as required by AGENTS.md. 
Do NOT use markdown TODO lists."
```

---

## 🔍 Problem: Beads Not Installed

### Symptom

Agent tries to use `bd` but command not found

### Diagnosis

```bash
# Check if bd is installed
which bd

# Try running bd
bd --help
```

### Solution

```bash
# Install Beads
curl -fsSL https://raw.githubusercontent.com/steveyegge/beads/main/scripts/install.sh | bash

# Or with npm
npm install -g @beadsio/beads

# Or with homebrew (macOS)
brew install beads

# Verify installation
bd --version
```

---

## 🔍 Problem: Beads Not Initialized

### Symptom

`bd: not in a beads repository`

### Diagnosis

```bash
# Check for .beads directory
ls -la .beads/

# Check if in git repo
git status
```

### Solution

```bash
# Must be in git repo first
git init
git config user.name "Your Name"
git config user.email "your.email@example.com"

# Then initialize beads
bd init

# Install git hooks (optional but recommended)
bd hooks install
```

---

## 🔍 Problem: Azure DevOps CLI Not Working

### Symptom

Agent tries to create work items but `az` command fails

### Diagnosis

```bash
# Check if Azure CLI installed
which az

# Check if logged in
az account show

# Check DevOps extension
az extension list | grep azure-devops
```

### Solution

```bash
# Install Azure CLI (Ubuntu/Debian)
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash

# Or with homebrew (macOS)
brew update && brew install azure-cli

# Install DevOps extension
az extension add --name azure-devops

# Login
az login

# Configure defaults
az devops configure --defaults \
  organization=https://dev.azure.com/YourOrg \
  project=YourProject
```

---

## 🔍 Problem: Agent Not Creating Azure DevOps Work Items

### Symptom

Beads created but no Azure DevOps work items

### Diagnosis

```bash
# Ask agent: "When should you create Azure DevOps work items?"
# Expected: Should reference decision table from beads-integration.md

# Check bead metadata
bd show bd-xxx --json | jq '.metadata.az_work_item'
# Should be null if no work item created
```

### Solutions

**1. Agent Doesn't Know When to Create Work Items**

```
# Add to prompt:
"This is a user-facing feature that needs Azure DevOps tracking. 
Please create both a bead and an Azure DevOps work item."
```

**2. Azure CLI Not Configured** See "Azure DevOps CLI Not Working" section above

**3. Agent Needs Explicit Instruction**

```
# Be specific in prompt:
"This is a P0 bug reported by users. Create:
1. A bead with -t bug -p 0
2. An Azure DevOps Bug work item
3. Link them bidirectionally"
```

---

## 🔍 Problem: Agent Creating Too Many Work Items

### Symptom

Internal tasks creating Azure DevOps work items

### Diagnosis

```bash
# Check which beads have AZ links
bd list --json | jq '.[] | select(.metadata.az_work_item != null) | {title, type, priority}'

# Verify they're user-facing
```

### Solution

```
# Add to prompt:
"This is internal refactoring/tech debt. 
Create a bead for tracking but DO NOT create an Azure DevOps work item."
```

**Decision Table Reference:**

- ✅ ALWAYS create AZ work items: User-facing features, bugs, epics, sprint
  commitments
- ⚠️ OPTIONAL: Tech debt (team discretion)
- ❌ NEVER: Internal planning, agent notes, duplicates

---

## 🔍 Problem: Beads Not Syncing

### Symptom

`bd sync` fails or beads lost between sessions

### Diagnosis

```bash
# Check git status
git status

# Check .beads directory
ls -la .beads/

# Check git log for bead commits
git log --oneline | grep bead
```

### Solutions

**1. Git Hooks Not Installed**

```bash
bd hooks install
git add .beads/
git commit -m "Add beads"
```

**2. .beads Directory Not Committed**

```bash
git add .beads/
git commit -m "Add beads tracking"
bd sync
```

**3. Manual Sync Required**

```bash
# Sync manually
bd sync

# Check sync worked
git log --oneline -1
```

---

## 🔍 Problem: Multi-Agent Conflicts

### Symptom

Multiple agents cause git conflicts in `.beads/`

### Diagnosis

```bash
# Check for conflicts
git status | grep conflict

# Look at .beads directory
cat .beads/issues.jsonl
```

### Solutions

**Should NOT happen** - Beads uses hash-based IDs to prevent conflicts!

**If it does happen:**

```bash
# 1. Check Beads version (might be bug)
bd --version

# 2. Each agent should sync before starting
bd sync  # Pull latest
bd ready # Check what's ready

# 3. Use separate branches per agent
git checkout -b agent-1-work
git checkout -b agent-2-work

# 4. Merge branches separately
git checkout main
git merge agent-1-work
git merge agent-2-work
```

**Prevention:**

- Each agent runs `bd sync` at start and end of session
- Each agent works on different tasks (check with `bd ready`)
- Use git branches for isolation

---

## 🔍 Problem: Context Lost Between Sessions

### Symptom

Agent doesn't remember previous tasks

### Diagnosis

```bash
# Check beads exist
bd list --json | jq '. | length'

# Check in-progress tasks
bd list --status in_progress
```

### Solutions

**1. Agent Not Checking Ready Tasks**

```
# Add to session start:
"Check bd ready to see what tasks you were working on"
```

**2. Beads Not Synced**

```bash
# Verify beads in git
git log --oneline | grep bead

# Sync if needed
bd sync
```

**3. Agent Not Querying Beads**

```
# Remind agent:
"Use bd list --status in_progress to find your previous work"
```

---

## 🔍 Problem: Dependencies Not Working

### Symptom

Agent works on blocked tasks

### Diagnosis

```bash
# Check dependencies
bd show bd-xxx --json | jq '.blocks, .blocked_by'

# Check ready tasks (should exclude blocked)
bd ready --json
```

### Solution

**1. Dependencies Not Set**

```bash
# Add dependency (Task A blocks on Task B)
bd dep add bd-taskA bd-taskB --type blocks

# Verify
bd show bd-taskA --json | jq '.blocked_by'
```

**2. Agent Not Checking Ready**

```
# Remind agent:
"Use bd ready to see only tasks that are ready (not blocked)"
```

---

## 🆘 Emergency Reset

**If everything is broken and you need to start fresh:**

```bash
# 1. Backup current state
cp -r .beads .beads.backup
git branch backup-$(date +%Y%m%d)

# 2. Remove beads
rm -rf .beads/

# 3. Reinitialize
bd init

# 4. Restore from backup if needed
# (only do this if you had working beads)
# cp -r .beads.backup/.beads/issues.jsonl .beads/

# 5. Commit fresh start
git add .beads/
git commit -m "Reinitialize beads"
```

---

## 📊 Health Check Commands

**Run these regularly to check Beads health:**

```bash
# Basic health
bd list --json | jq '. | length'  # Total beads
bd ready --json | jq '. | length'  # Ready to work

# Status breakdown
bd list --json | jq 'group_by(.status) | map({status: .[0].status, count: length})'

# Check for stale in-progress tasks
bd list --status in_progress --json

# Check Azure DevOps integration
bd list --json | jq '[.[] | select(.metadata.az_work_item != null)] | length'

# Check git sync
git log --oneline | head -5 | grep bead
```

---

## 🎯 Quick Fixes Summary

| Problem               | Quick Fix                            |
| --------------------- | ------------------------------------ |
| Agent not using beads | Add "Use Beads (bd)" to prompt       |
| `bd` not found        | Install: `curl ... \| bash`          |
| Not in beads repo     | `bd init`                            |
| Azure CLI fails       | `az login && az devops configure`    |
| No AZ work items      | Add "user-facing" to prompt          |
| Too many work items   | Add "internal only" to prompt        |
| Beads not syncing     | `bd sync` manually                   |
| Multi-agent conflicts | Use separate branches                |
| Lost context          | Check `bd list --status in_progress` |
| Blocked tasks         | Set dependencies with `bd dep add`   |

---

## 📞 Getting Help

**Documentation:**

- [beads-integration.md](assets/resources/processes/beads-integration.md) -
  Complete guide
- [TESTING_BEADS_WORKFLOW.md](TESTING_BEADS_WORKFLOW.md) - Test suite
- [AGENTS.md](AGENTS.md) - Required workflow

**External:**

- [Beads GitHub Issues](https://github.com/steveyegge/beads/issues)
- [Beads Documentation](https://github.com/steveyegge/beads/tree/main/docs)

**Still stuck?**

1. Run health check commands above
2. Check AGENTS.md is loaded in MCP server
3. Verify beads-integration.md is accessible
4. Test with simple prompt (see QUICK_TEST.md)

---

**Last Updated:** December 23, 2025\
**Status:** Ready for troubleshooting\
**Tip:** Bookmark this file for quick reference!
