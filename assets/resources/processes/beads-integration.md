# Beads Integration Guide for Azure DevOps Teams

**Status:** Required for all project work\
**Version:** 1.0\
**Last Updated:** December 2025

---

## Table of Contents

1. [Overview](#overview)
2. [Why Beads?](#why-beads)
3. [Installation & Setup](#installation--setup)
4. [Core Workflow](#core-workflow)
5. [Beads + Azure DevOps Integration](#beads--azure-devops-integration)
6. [Agent Instructions](#agent-instructions)
7. [Common Patterns](#common-patterns)
8. [Best Practices](#best-practices)
9. [Troubleshooting](#troubleshooting)
10. [Reference](#reference)

---

## Overview

**Beads (`bd`)** is a distributed, git-backed graph issue tracker designed
specifically for AI coding agents. It provides persistent, structured memory
that enables agents to handle long-horizon tasks without losing context.

### Key Concepts

- **Beads are git-based**: Issues stored as JSONL in `.beads/issues.jsonl`,
  versioned like code
- **Hash-based IDs**: Format `bd-a1b2` prevents merge conflicts in multi-agent
  workflows
- **Dependency graph**: Tasks can block/depend on other tasks
- **Zero-conflict merges**: Perfect for parallel agent work
- **Agent-optimized**: JSON output, auto-sync, background daemon

### Integration Philosophy

**Beads = Agent Memory | Azure DevOps = Team Visibility**

```
┌──────────────┐         Sync          ┌──────────────────┐
│              │ ◄──────────────────► │                  │
│    Beads     │                       │  Azure DevOps    │
│  (Agent DB)  │                       │   (Human UI)     │
│              │                       │                  │
└──────────────┘                       └──────────────────┘
     ↑                                         ↑
     │ Fast, structured                        │ Reports, sprints,
     │ queries for agents                      │ team collaboration
```

**When to create what:**

| Scenario                 | Create Bead? | Create AZ Work Item?               |
| ------------------------ | ------------ | ---------------------------------- |
| Agent breaks down work   | ✅ Always    | ✅ For human-facing tasks          |
| Agent detects bug/issue  | ✅ Always    | ✅ Always (link to bead)           |
| Agent planning sub-tasks | ✅ Always    | ⚠️ Optional (depends on scope)     |
| Human creates work item  | ⚠️ Optional  | ✅ Always                          |
| Multi-agent coordination | ✅ Required  | ✅ Required (one per epic/feature) |

---

## Why Beads?

### Problems Beads Solves

**❌ Without Beads (Traditional Approach):**

- Agent plans stored in markdown files (gets stale, lost between sessions)
- No structured task dependencies
- Merge conflicts in multi-agent workflows
- Context loss between sessions
- No audit trail of agent decisions

**✅ With Beads:**

- **Persistent memory**: Git-backed, survives sessions, branches, merges
- **Dependency-aware**: Agents know what tasks are ready, what's blocked
- **Zero conflicts**: Hash IDs prevent merge collisions
- **Fast queries**: `bd ready` shows immediately actionable work
- **Audit trail**: Every change tracked, who did what when
- **Multi-agent safe**: Multiple agents can work in parallel

### Integration Benefits

**For Agents:**

- Fast, structured queries (`bd ready --json`)
- Automatic context preservation across sessions
- Dependency tracking (know when tasks are unblocked)
- Git-backed versioning (every state is recoverable)

**For Teams:**

- Azure DevOps remains source of truth for human workflows
- Beads provides agent-optimized layer
- Bidirectional sync keeps systems in harmony
- Full audit trail in both systems

**For You:**

- Agents can autonomously manage complex, multi-session work
- Azure DevOps stays clean (no agent spam)
- Clear separation: beads for granular agent tasks, AZ for human-level features

---

## Installation & Setup

### 1. Install Beads

**macOS/Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/steveyegge/beads/main/scripts/install.sh | bash
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/steveyegge/beads/main/install.ps1 | iex
```

**npm:**

```bash
npm install -g @beads/bd
```

**Homebrew:**

```bash
brew install steveyegge/beads/bd
```

**Verify installation:**

```bash
bd version
```

### 2. Initialize in Your Repository

**Standard setup (commits to main branch):**

```bash
cd /path/to/your/repo
bd init
```

**Protected branch setup (commits to separate branch):**

```bash
# For repos with branch protection (recommended)
bd init --branch beads-metadata
```

This creates:

- `.beads/` directory with `issues.jsonl` and SQLite cache
- `.beads/` is git-tracked (team collaboration)
- Background daemon for auto-sync

**Note:** Use `--stealth` mode for personal use without committing to repo:

```bash
bd init --stealth
```

### 3. Install Git Hooks (Recommended)

Automatic sync on git operations:

```bash
bd hooks install
```

This installs hooks that:

- **pre-commit**: Flush pending changes before commit
- **post-merge**: Import updated beads after pull
- **pre-push**: Ensure JSONL is up-to-date before push
- **post-checkout**: Sync after branch checkout

### 4. Tell Your Agent

Add to your `AGENTS.md` or project instructions:

```markdown
# Task Tracking

**REQUIRED:** Use Beads (`bd`) for all task tracking and memory management.

- Create beads for work breakdown: `bd create "Task description" -p <priority>`
- Check ready work: `bd ready --json`
- Update progress: `bd update <id> --status in_progress`
- Close completed work: `bd close <id> --reason "Completed"`
- Show dependencies: `bd graph`

**Integration:** When creating beads that represent user-facing features or
bugs, also create corresponding Azure DevOps work items and link them.

See `assets/resources/processes/beads-integration.md` for full workflow.
```

---

## Core Workflow

### Basic Commands

```bash
# Create a task
bd create "Implement user authentication" -p 1

# List all tasks
bd list

# Show ready tasks (no blockers)
bd ready

# Show task details
bd show bd-a1b2

# Update task
bd update bd-a1b2 --status in_progress

# Close task
bd close bd-a1b2 --reason "Completed"

# Add dependency
bd dep add bd-child bd-parent  # child blocks on parent
```

### Hierarchical Tasks (Epics)

```bash
# Create epic
bd create "User Management Epic" -t epic -p 0
# Returns: bd-a3f8

# Create sub-tasks
bd create "Implement login" -p bd-a3f8 -t task
# Returns: bd-a3f8.1

bd create "Add unit tests for login" -p bd-a3f8.1 -t task
# Returns: bd-a3f8.1.1
```

**Hierarchy:**

- `bd-a3f8` - Epic
- `bd-a3f8.1` - Task (child of epic)
- `bd-a3f8.1.1` - Sub-task (child of task)

### Dependencies

```bash
# Task B blocks on Task A (A must complete before B)
bd dep add bd-B bd-A --type blocks

# Related tasks (informational link)
bd dep add bd-X bd-Y --type related

# Parent-child (structural hierarchy)
bd dep add bd-child bd-parent --type parent
```

### JSON Output (for Agents)

All commands support `--json` flag:

```bash
bd ready --json
bd show bd-a1b2 --json
bd create "New task" -p 1 --json
```

Returns structured JSON for easy parsing.

---

## Beads + Azure DevOps Integration

### Mapping Strategy

| Beads                      | Azure DevOps          | Notes                 |
| -------------------------- | --------------------- | --------------------- |
| `bd create` with `-t epic` | Epic                  | High-level features   |
| `bd create` with `-t task` | Task / User Story     | Mid-level work items  |
| `bd create` with `-t bug`  | Bug                   | Issue tracking        |
| `bd-X.Y.Z` hierarchy       | Parent-child links    | Epic → Feature → Task |
| `bd dep add --type blocks` | Predecessor-Successor | Dependency tracking   |
| `-p 0` (P0)                | Priority 1            | Critical              |
| `-p 1` (P1)                | Priority 2            | High                  |
| `-p 2` (P2)                | Priority 3            | Medium                |
| `-p 3` (P3)                | Priority 4            | Low                   |

### When to Create Azure DevOps Work Items

**✅ ALWAYS create AZ work item for:**

1. **User-facing features** - Anything visible to users or stakeholders
2. **Bugs** - All bugs should be tracked in Azure DevOps
3. **Epics/Stories** - High-level work that spans multiple sessions
4. **Sprint commitments** - Anything committed to a sprint

**⚠️ OPTIONAL for:**

1. **Agent planning tasks** - Internal breakdown (e.g., "Research API options")
2. **Technical debt** - Can stay in beads unless human visibility needed
3. **Sub-tasks** - Agent-level granularity (e.g., "Add unit test for X")

**❌ DO NOT create AZ work item for:**

1. **Transient notes** - Temporary agent memory
2. **Duplicate tracking** - Don't create both if bead is sufficient

### Integration Patterns

#### Pattern 1: Bead → Azure DevOps (Agent-Driven)

**Use case:** Agent discovers work during development

```bash
# 1. Agent creates bead
bd create "Fix validation bug in UserController" -t bug -p 0 --json

# 2. Get bead ID from JSON response
BEAD_ID="bd-abc123"

# 3. Create Azure DevOps work item
az boards work-item create \
  --type Bug \
  --title "Fix validation bug in UserController" \
  --description "Linked to bead: $BEAD_ID" \
  --assigned-to "me@example.com" \
  --project "MyProject"

# 4. Link in bead metadata (optional)
bd update $BEAD_ID --metadata "az_work_item=12345"
```

#### Pattern 2: Azure DevOps → Bead (Human-Driven)

**Use case:** Human creates work item, agent needs to track it

```bash
# 1. Human creates work item in Azure DevOps UI
# Work Item ID: 12345

# 2. Agent imports to beads
bd create "Implement feature X (AZ-12345)" -p 1 --json

# 3. Link in metadata
bd update bd-xyz --metadata "az_work_item=12345"
```

#### Pattern 3: Bidirectional Sync (Automation)

**Use case:** Keep systems in sync automatically

**Example: Create AZ work item when bead is created**

```bash
# .beads-hooks/post-create.sh
#!/bin/bash
# Triggered after bd create

BEAD_ID="$1"
BEAD_TITLE="$2"
BEAD_TYPE="$3"
BEAD_PRIORITY="$4"

# Only sync certain bead types
if [[ "$BEAD_TYPE" == "epic" || "$BEAD_TYPE" == "bug" ]]; then
  # Map bead priority to AZ priority
  case "$BEAD_PRIORITY" in
    0) AZ_PRIORITY=1 ;;
    1) AZ_PRIORITY=2 ;;
    2) AZ_PRIORITY=3 ;;
    *) AZ_PRIORITY=4 ;;
  esac

  # Create AZ work item
  AZ_TYPE="Task"
  [[ "$BEAD_TYPE" == "bug" ]] && AZ_TYPE="Bug"
  [[ "$BEAD_TYPE" == "epic" ]] && AZ_TYPE="Epic"

  az boards work-item create \
    --type "$AZ_TYPE" \
    --title "$BEAD_TITLE" \
    --description "Linked bead: $BEAD_ID" \
    --priority "$AZ_PRIORITY" \
    --project "MyProject" \
    --output json > /tmp/az_work_item.json

  # Extract work item ID and link back
  AZ_ID=$(jq -r '.id' /tmp/az_work_item.json)
  bd update "$BEAD_ID" --metadata "az_work_item=$AZ_ID"

  echo "Created AZ work item $AZ_ID for bead $BEAD_ID"
fi
```

**Install the hook:**

```bash
chmod +x .beads-hooks/post-create.sh
```

**Beads will automatically execute hooks in `.beads-hooks/` directory.**

---

## Agent Instructions

### For AI Coding Agents

**REQUIRED WORKFLOW:**

1. **Start every session** by checking ready work:
   ```bash
   bd ready --json
   ```

2. **Create beads for all work** (never use markdown TODO lists):
   ```bash
   bd create "Implement feature X" -p 1 --json
   bd create "Add tests for feature X" -p 1 --json
   bd create "Update documentation" -p 2 --json
   ```

3. **Create Azure DevOps work items** for user-facing work:
   ```bash
   # If bead represents a bug, feature, or epic
   az boards work-item create --type <type> --title "..." --project "MyProject"
   ```

4. **Update status** as you work:
   ```bash
   bd update bd-abc --status in_progress --json
   ```

5. **Close completed work**:
   ```bash
   bd close bd-abc --reason "Completed" --json
   ```

6. **End of session sync**:
   ```bash
   bd sync
   ```

### Agent Session Template

```bash
# === START OF SESSION ===

# 1. Check what's ready
bd ready --json

# 2. Choose task
bd show bd-abc --json

# 3. Update status
bd update bd-abc --status in_progress

# === DO THE WORK ===

# 4. Create sub-tasks as needed
bd create "Sub-task for bd-abc" -p bd-abc --json

# 5. Create AZ work item if needed
if [[ "user-facing work" ]]; then
  az boards work-item create --type Task --title "..." --project "MyProject"
  bd update bd-abc --metadata "az_work_item=12345"
fi

# === END OF SESSION ===

# 6. Close completed work
bd close bd-abc --reason "Completed"

# 7. Sync to git
bd sync

# 8. Summary for human
bd list --status open --json
```

### Multi-Agent Coordination

**Scenario:** Multiple agents working in parallel

```bash
# Agent 1: Feature development
bd create "Implement login API" -p 0 --json
bd update bd-login --status in_progress

# Agent 2: Testing (depends on Agent 1)
bd create "Test login API" -p 1 --json
bd dep add bd-test bd-login --type blocks

# Agent 2 checks ready work
bd ready --json
# Will NOT show bd-test (blocked by bd-login)

# Agent 1 completes work
bd close bd-login --reason "Completed"
bd sync

# Agent 2 pulls changes
git pull
bd ready --json
# NOW shows bd-test (unblocked)
```

---

## Common Patterns

### Pattern: Epic Breakdown

```bash
# 1. Create epic
bd create "User Authentication System" -t epic -p 0
# Returns: bd-auth

# 2. Break down into features
bd create "Login API endpoint" -p bd-auth -t task
# Returns: bd-auth.1

bd create "JWT token generation" -p bd-auth -t task
# Returns: bd-auth.2

bd create "Password reset flow" -p bd-auth -t task
# Returns: bd-auth.3

# 3. Further breakdown
bd create "Write login controller" -p bd-auth.1 -t task
# Returns: bd-auth.1.1

bd create "Add authentication middleware" -p bd-auth.1 -t task
# Returns: bd-auth.1.2

# 4. Add dependencies
bd dep add bd-auth.1.2 bd-auth.1.1 --type blocks

# 5. Check what's ready
bd ready --json
# Returns: bd-auth.1.1 (no blockers)

# 6. Create Azure DevOps Epic for visibility
az boards work-item create \
  --type Epic \
  --title "User Authentication System" \
  --description "Bead: bd-auth" \
  --project "MyProject"
```

### Pattern: Bug Workflow

```bash
# 1. Agent discovers bug
bd create "NullReferenceException in UserController.Login" -t bug -p 0 --json

# 2. Create Azure DevOps bug
az boards work-item create \
  --type Bug \
  --title "NullReferenceException in UserController.Login" \
  --description "Bead: bd-bug123" \
  --assigned-to "developer@example.com" \
  --project "MyProject"

# 3. Link them
bd update bd-bug123 --metadata "az_work_item=54321"

# 4. Work on fix
bd update bd-bug123 --status in_progress

# 5. Create test task
bd create "Add regression test for bug bd-bug123" -p bd-bug123 -t task

# 6. Close bug
bd close bd-bug123 --reason "Fixed in PR #123"

# 7. Update Azure DevOps (manual or automated)
az boards work-item update 54321 --state Resolved
```

### Pattern: Sprint Planning

```bash
# 1. List open work
bd list --status open --json

# 2. Prioritize for sprint
bd update bd-abc --priority 0  # P0 = Critical
bd update bd-def --priority 1  # P1 = High
bd update bd-ghi --priority 2  # P2 = Medium

# 3. Add sprint metadata
bd update bd-abc --metadata "sprint=2025-Q1-Sprint1"
bd update bd-def --metadata "sprint=2025-Q1-Sprint1"

# 4. Query sprint work
bd list --metadata "sprint=2025-Q1-Sprint1" --json

# 5. Create corresponding AZ work items for sprint
# (Manual or automated sync)
```

### Pattern: Multi-Agent Project

```bash
# === Project Lead Agent ===
bd create "Migrate to .NET 8" -t epic -p 0
# Returns: bd-migration

bd create "Update NuGet packages" -p bd-migration
bd create "Update C# language features" -p bd-migration
bd create "Update EF Core" -p bd-migration
bd create "Run regression tests" -p bd-migration

# Add dependencies
bd dep add bd-regression bd-efcore --type blocks
bd dep add bd-regression bd-nuget --type blocks
bd dep add bd-regression bd-csharp --type blocks

bd sync
git push

# === Agent 1 (NuGet) ===
git pull
bd ready --json  # Shows bd-nuget (no blockers)
bd update bd-nuget --status in_progress
# ... do work ...
bd close bd-nuget --reason "Completed"
bd sync
git push

# === Agent 2 (C#) ===
git pull
bd ready --json  # Shows bd-csharp
bd update bd-csharp --status in_progress
# ... do work ...
bd close bd-csharp --reason "Completed"
bd sync
git push

# === Agent 3 (Testing) ===
git pull
bd ready --json  # NOW shows bd-regression (unblocked)
bd update bd-regression --status in_progress
# ... do work ...
bd close bd-regression --reason "All tests pass"
bd sync
git push
```

---

## Best Practices

### For Agents

1. **Always use beads instead of markdown plans**
   - ❌ Bad: `TODO: Implement feature X`
   - ✅ Good: `bd create "Implement feature X" -p 1`

2. **Create beads at appropriate granularity**
   - Epics: Large features (> 1 week)
   - Tasks: Medium work (1-3 days)
   - Sub-tasks: Small work (< 1 day)

3. **Use dependencies to model blockers**
   - If Task B can't start until Task A is done, use `bd dep add B A`

4. **Always sync at end of session**
   ```bash
   bd sync
   ```

5. **Use JSON output for programmatic access**
   ```bash
   bd ready --json | jq '.[] | select(.priority == 0)'
   ```

6. **Link to Azure DevOps for visibility**
   - Add `az_work_item` metadata to beads
   - Reference bead IDs in AZ work item descriptions

7. **Close beads with meaningful reasons**
   ```bash
   bd close bd-abc --reason "Completed in PR #456"
   bd close bd-def --reason "Duplicate of bd-xyz"
   bd close bd-ghi --reason "Won't fix - out of scope"
   ```

### For Teams

1. **Document your integration approach**
   - Which beads get AZ work items?
   - Who creates what?
   - How are they linked?

2. **Use protected branch mode for main branches**
   ```bash
   bd init --branch beads-metadata
   ```

3. **Install git hooks for automatic sync**
   ```bash
   bd hooks install
   ```

4. **Establish naming conventions**
   - Prefix beads with context: `"API: Add endpoint"`, `"UI: Fix button"`
   - Reference AZ work items: `"Implement feature X (AZ-12345)"`

5. **Regular cleanup**
   ```bash
   # Close stale beads
   bd list --status open --json | jq -r '.[] | select(.created < "2024-01-01") | .id'

   # Compact old closed beads (optional)
   bd compact --days 90
   ```

### Integration Best Practices

1. **Bidirectional linking**
   - Beads reference AZ work item IDs in metadata
   - AZ work items reference bead IDs in description

2. **Single source of truth**
   - Beads: Agent memory, task dependencies, granular breakdown
   - Azure DevOps: Team visibility, sprint planning, reporting

3. **Automation over manual sync**
   - Use `.beads-hooks/` for automatic AZ work item creation
   - Avoid manual copy-paste

4. **Clear ownership**
   - Agent-created beads → Agent owns
   - Human-created AZ work items → Human owns
   - Sync when crossing boundaries

---

## Troubleshooting

### Common Issues

**Issue: Merge conflicts in `.beads/issues.jsonl`**

```bash
# Accept remote version and re-import
git checkout --theirs .beads/issues.jsonl
bd import -i .beads/issues.jsonl

# Or accept local version
git checkout --ours .beads/issues.jsonl
bd import -i .beads/issues.jsonl
```

**Issue: Database out of sync with JSONL**

```bash
# Force export and re-import
bd sync
```

**Issue: Daemon not running**

```bash
# Check daemon status
bd daemon status

# Start daemon
bd daemon start

# Restart daemon
bd daemon restart
```

**Issue: Git hooks not working**

```bash
# Reinstall hooks
bd hooks install

# Verify hooks are executable
ls -la .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

**Issue: Can't find ready work**

```bash
# Check if tasks are blocked
bd list --json | jq '.[] | select(.status == "open") | {id, title, blockers}'

# Show dependency graph
bd graph
```

### Debugging Commands

```bash
# Show all beads with details
bd list --json | jq

# Show specific bead
bd show bd-abc --json

# Show dependencies
bd deps bd-abc

# Show reverse dependencies (what blocks on this)
bd rdeps bd-abc

# Export database to JSONL
bd export -o /tmp/beads-backup.jsonl

# Import from JSONL
bd import -i /tmp/beads-backup.jsonl

# Check database integrity
bd doctor
```

---

## Reference

### Bead Properties

```json
{
  "id": "bd-a1b2c3",
  "title": "Implement feature X",
  "type": "task",
  "status": "open",
  "priority": 1,
  "created": "2025-12-23T10:00:00Z",
  "updated": "2025-12-23T11:30:00Z",
  "closed": null,
  "assignee": "agent-1",
  "description": "Detailed description",
  "metadata": {
    "az_work_item": "12345",
    "sprint": "2025-Q1-Sprint1"
  },
  "parent": "bd-parent",
  "children": ["bd-child1", "bd-child2"],
  "blockers": ["bd-blocker"],
  "blocked_by": ["bd-task1"]
}
```

### Type Values

- `epic` - High-level feature (multiple weeks)
- `task` - Mid-level work (1-3 days)
- `bug` - Issue to fix
- `chore` - Maintenance work
- `research` - Investigation task

### Status Values

- `open` - Not started
- `in_progress` - Currently being worked on
- `blocked` - Waiting on dependency
- `closed` - Completed

### Priority Values

- `0` (P0) - Critical (drop everything)
- `1` (P1) - High (do soon)
- `2` (P2) - Medium (normal priority)
- `3` (P3) - Low (nice to have)

### Dependency Types

- `blocks` - Parent must complete before child can start
- `related` - Informational link (no blocking)
- `parent` - Structural hierarchy (epic → task → sub-task)

### Useful Commands

```bash
# Create
bd create "Title" -p <priority> -t <type> --parent <parent-id> --json

# List
bd list --status <status> --priority <priority> --type <type> --json

# Ready (no blockers)
bd ready --json

# Show
bd show <id> --json

# Update
bd update <id> --status <status> --priority <priority> --json

# Close
bd close <id> --reason "message" --json

# Dependencies
bd dep add <child> <parent> --type <type>
bd dep remove <child> <parent>
bd deps <id>
bd rdeps <id>

# Graph
bd graph
bd graph <id>  # Sub-graph for specific bead

# Sync
bd sync

# Daemon
bd daemon start
bd daemon stop
bd daemon restart
bd daemon status

# Hooks
bd hooks install
bd hooks uninstall

# Export/Import
bd export -o <file>
bd import -i <file>

# Doctor (check integrity)
bd doctor
```

---

## Summary

**Beads + Azure DevOps = Optimal Workflow**

- **Beads** provides agent memory, dependency tracking, zero-conflict merges
- **Azure DevOps** provides team visibility, sprint planning, reporting
- **Integration** keeps both systems synchronized and serving their strengths

**Key Takeaway:** Whenever an agent creates a bead that represents user-facing
work (epics, bugs, features), also create an Azure DevOps work item and link
them bidirectionally.

**Next Steps:**

1. Install Beads:
   `curl -fsSL https://raw.githubusercontent.com/steveyegge/beads/main/scripts/install.sh | bash`
2. Initialize: `bd init`
3. Install hooks: `bd hooks install`
4. Update `AGENTS.md` to require beads usage
5. Set up automation hooks in `.beads-hooks/`

For more information:

- Beads documentation: https://github.com/steveyegge/beads
- Azure DevOps CLI: https://docs.microsoft.com/en-us/cli/azure/boards
- Team standards: See other guides in `assets/resources/`
