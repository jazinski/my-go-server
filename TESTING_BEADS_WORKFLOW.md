# Testing Beads Workflow Compliance

**Purpose:** Verify that AI agents correctly follow the required Beads workflow
when working on projects.

---

## 🎯 What We're Testing

1. **Agents use Beads for task tracking** (not markdown TODO lists)
2. **Agents follow the required workflow** (ready → create → update → close →
   sync)
3. **Agents create Azure DevOps work items** when appropriate
4. **Agents link beads ↔ Azure DevOps** bidirectionally
5. **Agents coordinate properly** in multi-agent scenarios

---

## 🧪 Test Suite

### Test 1: Basic Workflow Compliance

**Objective:** Verify agent uses Beads for task tracking

**Setup:**

```bash
# Create test project
mkdir /tmp/beads-test
cd /tmp/beads-test
git init
bd init
```

**Test Prompt:**

```
"I need you to add a new feature: user authentication with OAuth2. 
Please create tasks for this work using our standard workflow."
```

**Expected Behavior:** ✅ Agent runs `bd create` commands (NOT creates markdown
TODO list) ✅ Agent creates multiple beads for sub-tasks ✅ Agent uses proper
priorities (-p flag) ✅ Agent uses proper types (-t flag: epic, task, bug)

**Verification:**

```bash
# Check beads were created
bd list --json

# Should see beads with structure like:
# - Epic: "Add OAuth2 authentication" (bd-xxxx)
#   - Task: "Implement OAuth2 provider configuration" (bd-xxxx.1)
#   - Task: "Add login/logout endpoints" (bd-xxxx.2)
#   - Task: "Add user session management" (bd-xxxx.3)
```

**Pass Criteria:**

- ✅ Agent uses `bd create` commands
- ✅ Creates hierarchical tasks (epic → tasks)
- ✅ No markdown TODO lists created
- ✅ All beads have proper metadata

---

### Test 2: Session Workflow

**Objective:** Verify agent follows start/end session workflow

**Test Prompt:**

```
"Start working on implementing the OAuth2 feature we discussed."
```

**Expected Behavior:** ✅ Agent starts with `bd ready --json` to check ready
tasks ✅ Agent selects an appropriate ready task ✅ Agent runs
`bd update bd-xxx --status in_progress` ✅ Agent does the work ✅ Agent runs
`bd close bd-xxx --reason "Completed"` ✅ Agent runs `bd sync` at end

**Verification:**

```bash
# Check task status progression
bd show bd-xxx --json | jq '.status'
# Should show: open → in_progress → closed

# Check git commits
git log --oneline
# Should see bead sync commits
```

**Pass Criteria:**

- ✅ Agent checks `bd ready` first
- ✅ Updates status to `in_progress` before work
- ✅ Closes bead after completion
- ✅ Syncs at end of session

---

### Test 3: Azure DevOps Integration

**Objective:** Verify agent creates Azure DevOps work items when appropriate

**Setup:**

```bash
# Configure Azure DevOps (use your real project)
az devops configure --defaults \
  organization=https://dev.azure.com/YourOrg \
  project=YourTestProject
```

**Test Prompt:**

```
"Create a task to fix the login validation bug that users are reporting. 
This is a P0 bug that needs to be tracked in Azure DevOps."
```

**Expected Behavior:** ✅ Agent creates bead:
`bd create "Fix login validation bug" -t bug -p 0` ✅ Agent recognizes this is
user-facing and P0 ✅ Agent creates Azure DevOps work item:
`az boards work-item create` ✅ Agent links them bidirectionally

**Verification:**

```bash
# Check bead was created
bd list --json | jq '.[] | select(.title | contains("login validation"))'

# Check metadata has AZ work item ID
bd show bd-xxx --json | jq '.metadata.az_work_item'
# Should return work item number like "12345"

# Check Azure DevOps work item exists
az boards work-item show --id 12345
# Should show work item with bead ID in description
```

**Pass Criteria:**

- ✅ Bead created with correct type (bug) and priority (P0)
- ✅ Azure DevOps work item created
- ✅ Bead metadata contains `az_work_item` ID
- ✅ AZ work item description contains bead ID

---

### Test 4: When NOT to Create Azure DevOps Work Items

**Objective:** Verify agent doesn't create unnecessary AZ work items

**Test Prompt:**

```
"I need you to refactor the internal caching logic to improve performance. 
This is technical debt that doesn't need external tracking."
```

**Expected Behavior:** ✅ Agent creates bead for the work ✅ Agent does NOT
create Azure DevOps work item (internal tech debt) ✅ Agent explains reasoning

**Verification:**

```bash
# Check bead exists
bd list --json | jq '.[] | select(.title | contains("caching"))'

# Check NO Azure DevOps link
bd show bd-xxx --json | jq '.metadata.az_work_item'
# Should return null or not exist
```

**Pass Criteria:**

- ✅ Bead created for tracking
- ✅ NO Azure DevOps work item created
- ✅ Agent explains it's internal work

---

### Test 5: Multi-Agent Coordination

**Objective:** Verify multiple agents can work without conflicts

**Setup:**

```bash
# Create epic with sub-tasks
cd /tmp/beads-test
bd create "Implement user dashboard" -t epic -p 1
# Note the epic ID, e.g., bd-a3f8

bd create "Backend API for dashboard" -t task -p 1 --parent bd-a3f8
bd create "Frontend dashboard UI" -t task -p 1 --parent bd-a3f8
bd create "Dashboard tests" -t task -p 1 --parent bd-a3f8

bd sync
```

**Test Prompt (to multiple AI agents simultaneously):**

**Agent 1:**

```
"Work on the backend API for the user dashboard (check ready tasks)"
```

**Agent 2:**

```
"Work on the frontend dashboard UI (check ready tasks)"
```

**Expected Behavior:** ✅ Each agent runs `bd ready` independently ✅ Each agent
picks different tasks ✅ Each agent updates their task to `in_progress` ✅ No
conflicts when syncing

**Verification:**

```bash
# Check different tasks are in progress
bd list --json | jq '.[] | select(.status == "in_progress") | .title'
# Should show both tasks

# Simulate git merge (no conflicts expected)
git log --oneline --all
# Should see commits from both agents with different bead IDs
```

**Pass Criteria:**

- ✅ Agents select different tasks
- ✅ Both update status independently
- ✅ No git merge conflicts
- ✅ Each agent syncs successfully

---

### Test 6: Dependency Management

**Objective:** Verify agent respects task dependencies

**Setup:**

```bash
cd /tmp/beads-test

# Create tasks with dependencies
bd create "Design database schema" -t task -p 1
# Returns bd-b1c2

bd create "Implement data access layer" -t task -p 1
# Returns bd-d3e4

# Add dependency: DAL blocks on schema design
bd dep add bd-d3e4 bd-b1c2 --type blocks
```

**Test Prompt:**

```
"What tasks are ready to work on? Start working on the highest priority task."
```

**Expected Behavior:** ✅ Agent runs `bd ready` ✅ Agent sees only "Design
database schema" is ready ✅ Agent does NOT work on "Implement data access
layer" (blocked) ✅ Agent updates the ready task

**Verification:**

```bash
# Check ready tasks
bd ready --json
# Should only show bd-b1c2 (schema design)
# Should NOT show bd-d3e4 (DAL - blocked)

# Check dependency
bd show bd-d3e4 --json | jq '.blocks'
# Should show dependency on bd-b1c2
```

**Pass Criteria:**

- ✅ Agent only works on ready tasks
- ✅ Agent does not work on blocked tasks
- ✅ Agent respects dependency graph

---

### Test 7: Context Persistence

**Objective:** Verify agent maintains context across sessions

**Session 1 Prompt:**

```
"Create a task to implement email notifications and start working on it. 
Stop when you've implemented the basic structure."
```

**Expected Behavior (Session 1):** ✅ Agent creates bead ✅ Agent updates status
to `in_progress` ✅ Agent implements partial work ✅ Agent syncs but does NOT
close bead

**Exit and restart agent (simulate new session)**

**Session 2 Prompt:**

```
"Continue working on the task you were doing before."
```

**Expected Behavior (Session 2):** ✅ Agent runs `bd list --status in_progress`
or `bd ready` ✅ Agent finds the email notifications task ✅ Agent continues
from where it left off ✅ Agent closes bead when done

**Verification:**

```bash
# Check task exists and has history
bd show bd-xxx --json | jq '.history'
# Should show multiple status updates across sessions

# Check work was actually continued
git log --oneline -- email-notifications/
# Should see commits from both sessions
```

**Pass Criteria:**

- ✅ Agent finds previous in-progress task
- ✅ Agent continues work seamlessly
- ✅ Context maintained across sessions
- ✅ Work completed and bead closed

---

## 🔍 Automated Test Script

Here's a script to run all tests automatically:

```bash
#!/bin/bash
# test-beads-workflow.sh

set -e

echo "🧪 Testing Beads Workflow Compliance"
echo "===================================="
echo ""

# Setup
TEST_DIR="/tmp/beads-workflow-test-$(date +%s)"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

git init
git config user.name "Test User"
git config user.email "test@example.com"
bd init

echo "✅ Test environment created: $TEST_DIR"
echo ""

# Test 1: Basic workflow
echo "📋 Test 1: Basic Workflow Compliance"
echo "Expected: Agent should use 'bd create' commands"
echo "Prompt your agent: 'Create tasks for implementing user authentication'"
echo ""
read -p "Did agent use 'bd create' (not markdown TODO)? (y/n): " response
if [ "$response" = "y" ]; then
    echo "✅ PASS: Agent used Beads"
else
    echo "❌ FAIL: Agent did not use Beads"
    exit 1
fi
echo ""

# Test 2: Check beads exist
echo "📋 Test 2: Verify Beads Created"
BEAD_COUNT=$(bd list --json | jq '. | length')
echo "Found $BEAD_COUNT beads"
if [ "$BEAD_COUNT" -gt 0 ]; then
    echo "✅ PASS: Beads created"
    bd list
else
    echo "❌ FAIL: No beads found"
    exit 1
fi
echo ""

# Test 3: Session workflow
echo "📋 Test 3: Session Workflow"
echo "Prompt agent: 'Start working on the first task'"
echo ""
read -p "Did agent start with 'bd ready'? (y/n): " response
if [ "$response" = "y" ]; then
    echo "✅ PASS: Agent checked ready tasks"
else
    echo "❌ FAIL: Agent skipped bd ready"
    exit 1
fi
echo ""

read -p "Did agent update status to 'in_progress'? (y/n): " response
if [ "$response" = "y" ]; then
    echo "✅ PASS: Agent updated status"
else
    echo "❌ FAIL: Agent did not update status"
    exit 1
fi
echo ""

# Test 4: Azure DevOps integration
echo "📋 Test 4: Azure DevOps Integration"
echo "Prompt agent: 'Create a P0 bug for login validation (user-facing)'"
echo ""
read -p "Did agent create Azure DevOps work item? (y/n): " response
if [ "$response" = "y" ]; then
    echo "✅ PASS: Agent created AZ work item"
    read -p "Enter bead ID to verify: " bead_id
    if bd show "$bead_id" --json | jq -e '.metadata.az_work_item' > /dev/null; then
        echo "✅ PASS: Bead linked to AZ work item"
    else
        echo "❌ FAIL: Bead not linked to AZ work item"
        exit 1
    fi
else
    echo "❌ FAIL: Agent did not create AZ work item"
    exit 1
fi
echo ""

# Test 5: When NOT to create AZ work items
echo "📋 Test 5: No AZ Work Item for Internal Tasks"
echo "Prompt agent: 'Refactor internal caching (tech debt)'"
echo ""
read -p "Did agent create bead but NOT AZ work item? (y/n): " response
if [ "$response" = "y" ]; then
    echo "✅ PASS: Agent correctly skipped AZ work item"
else
    echo "❌ FAIL: Agent created unnecessary AZ work item"
    exit 1
fi
echo ""

# Summary
echo "=================================="
echo "🎉 All Tests Passed!"
echo "=================================="
echo ""
echo "Test directory: $TEST_DIR"
echo "Review with: cd $TEST_DIR && bd list"
echo ""
```

**Save and run:**

```bash
chmod +x test-beads-workflow.sh
./test-beads-workflow.sh
```

---

## 📊 Compliance Checklist

Use this checklist during testing:

### Required Behaviors ✅

- [ ] Agent uses `bd` commands (not markdown TODO lists)
- [ ] Agent starts sessions with `bd ready --json`
- [ ] Agent creates beads for all work
- [ ] Agent updates status to `in_progress` before work
- [ ] Agent closes beads with `bd close` after completion
- [ ] Agent syncs with `bd sync` at end of session
- [ ] Agent creates Azure DevOps work items for user-facing features
- [ ] Agent creates Azure DevOps work items for bugs (P0, P1)
- [ ] Agent links beads ↔ Azure DevOps bidirectionally
- [ ] Agent does NOT create AZ work items for internal tasks
- [ ] Agent respects task dependencies
- [ ] Agent coordinates with other agents without conflicts
- [ ] Agent maintains context across sessions

### Prohibited Behaviors ❌

- [ ] Agent creates markdown TODO lists instead of beads
- [ ] Agent skips `bd ready` at session start
- [ ] Agent works on blocked tasks
- [ ] Agent forgets to update status
- [ ] Agent forgets to close completed beads
- [ ] Agent forgets to sync at end
- [ ] Agent creates duplicate AZ work items
- [ ] Agent creates AZ work items for internal work

---

## 🔧 Debugging Failed Tests

### Agent Not Using Beads

**Symptom:** Agent creates markdown TODO lists

**Diagnosis:**

```bash
# Check if agent has access to AGENTS.md
# Agent should see the required Beads workflow section
```

**Fix:**

- Ensure MCP server is running and connected
- Verify AGENTS.md is loaded as a resource
- Restart agent/OpenCode

### Agent Not Creating Azure DevOps Work Items

**Symptom:** No work items created for user-facing features

**Diagnosis:**

```bash
# Check Azure CLI is configured
az devops configure --list

# Check agent has access to beads-integration.md
# Should see decision table for when to create AZ work items
```

**Fix:**

- Configure Azure DevOps CLI: `az login` and `az devops configure`
- Ensure agent reads the integration guide
- Explicitly mention "This needs Azure DevOps tracking"

### Agent Creating Too Many Azure DevOps Work Items

**Symptom:** Work items created for internal tasks

**Diagnosis:**

- Agent not following decision table
- Agent not distinguishing user-facing vs internal work

**Fix:**

- Emphasize in prompt: "This is internal work"
- Review beads-integration.md decision table
- Update agent instructions if needed

### Context Not Persisting

**Symptom:** Agent doesn't remember previous tasks

**Diagnosis:**

```bash
# Check beads are synced to git
git log --oneline | grep bead

# Check .beads directory exists
ls -la .beads/
```

**Fix:**

- Ensure agent runs `bd sync` at end of session
- Check git hooks are installed: `bd hooks install`
- Verify .beads directory is committed

---

## 📈 Monitoring Agent Compliance

### Daily Checks

```bash
# Check how many beads created today
bd list --json | jq '.[] | select(.created_at | startswith("2025-12-23"))'

# Check beads without Azure DevOps links (should be mostly internal work)
bd list --json | jq '.[] | select(.metadata.az_work_item == null) | .title'

# Check open tasks
bd list --status open --json | jq '. | length'

# Check in-progress tasks (should be few)
bd list --status in_progress --json
```

### Weekly Audit

```bash
# Generate report
cat > weekly-beads-report.sh << 'EOF'
#!/bin/bash
echo "Weekly Beads Compliance Report"
echo "=============================="
echo ""
echo "Total Beads: $(bd list --json | jq '. | length')"
echo "Open: $(bd list --status open --json | jq '. | length')"
echo "In Progress: $(bd list --status in_progress --json | jq '. | length')"
echo "Closed: $(bd list --status closed --json | jq '. | length')"
echo ""
echo "With Azure DevOps Links: $(bd list --json | jq '[.[] | select(.metadata.az_work_item != null)] | length')"
echo "Without AZ Links: $(bd list --json | jq '[.[] | select(.metadata.az_work_item == null)] | length')"
echo ""
echo "By Type:"
echo "  Epics: $(bd list --json | jq '[.[] | select(.type == "epic")] | length')"
echo "  Tasks: $(bd list --json | jq '[.[] | select(.type == "task")] | length')"
echo "  Bugs: $(bd list --json | jq '[.[] | select(.type == "bug")] | length')"
EOF
chmod +x weekly-beads-report.sh
./weekly-beads-report.sh
```

---

## 🎯 Success Metrics

### Good Compliance

- ✅ 90%+ of work tracked in beads
- ✅ 0 markdown TODO lists created
- ✅ 100% of user-facing features have AZ work items
- ✅ 100% of P0/P1 bugs have AZ work items
- ✅ <20% of beads have AZ links (most work is internal)
- ✅ No git conflicts from beads
- ✅ All beads eventually closed (not abandoned)

### Red Flags

- ❌ Markdown TODO lists appearing
- ❌ Tasks not tracked in beads
- ❌ Many tasks stuck in `in_progress`
- ❌ Internal tasks creating AZ work items
- ❌ Git merge conflicts in .beads/
- ❌ Beads created but never updated

---

## 📚 Reference

**Related Documentation:**

- [beads-integration.md](../resources/processes/beads-integration.md) - Complete
  workflow guide
- [beads-azure-devops-automation.md](../resources/examples/beads-azure-devops-automation.md) -
  Automation scripts
- [AGENTS.md](../../AGENTS.md) - Required agent workflow

**External Resources:**

- [Beads GitHub](https://github.com/steveyegge/beads)
- [Beads Documentation](https://github.com/steveyegge/beads/blob/main/docs/)

---

**Last Updated:** December 23, 2025 **Status:** Ready for testing **Next
Action:** Run test suite with your AI agents
