# Beads + Azure DevOps Automation Examples

This document provides practical automation examples for integrating Beads issue
tracking with Azure DevOps work items.

---

## Table of Contents

1. [Automatic Work Item Creation](#automatic-work-item-creation)
2. [Status Synchronization](#status-synchronization)
3. [Agent Session Scripts](#agent-session-scripts)
4. [Multi-Agent Coordination](#multi-agent-coordination)
5. [PowerShell Examples](#powershell-examples)

---

## Automatic Work Item Creation

### Bash Hook: Auto-create Azure DevOps Work Items

Save this as `.beads-hooks/post-create.sh` in your repository:

```bash
#!/bin/bash
# .beads-hooks/post-create.sh
# Automatically creates Azure DevOps work items when beads are created

set -e

BEAD_ID="$1"
BEAD_TITLE="$2"
BEAD_TYPE="$3"
BEAD_PRIORITY="$4"

# Configuration
AZ_PROJECT="MyProject"
AZ_ASSIGNED_TO="me@example.com"

# Only sync certain bead types (bugs, epics, and P0/P1 tasks)
if [[ "$BEAD_TYPE" == "epic" || "$BEAD_TYPE" == "bug" ]] || \
   [[ "$BEAD_TYPE" == "task" && ( "$BEAD_PRIORITY" == "0" || "$BEAD_PRIORITY" == "1" ) ]]; then
  
  echo "Creating Azure DevOps work item for bead $BEAD_ID..."
  
  # Map bead priority to Azure DevOps priority
  case "$BEAD_PRIORITY" in
    0) AZ_PRIORITY=1 ;;  # P0 = Critical
    1) AZ_PRIORITY=2 ;;  # P1 = High
    2) AZ_PRIORITY=3 ;;  # P2 = Medium
    *) AZ_PRIORITY=4 ;;  # P3 = Low
  esac
  
  # Map bead type to Azure DevOps work item type
  AZ_TYPE="Task"
  [[ "$BEAD_TYPE" == "bug" ]] && AZ_TYPE="Bug"
  [[ "$BEAD_TYPE" == "epic" ]] && AZ_TYPE="Epic"
  
  # Create work item
  az boards work-item create \
    --type "$AZ_TYPE" \
    --title "$BEAD_TITLE" \
    --description "Linked bead: $BEAD_ID\n\nAutomatically created from Beads." \
    --priority "$AZ_PRIORITY" \
    --assigned-to "$AZ_ASSIGNED_TO" \
    --project "$AZ_PROJECT" \
    --output json > /tmp/az_work_item_$BEAD_ID.json
  
  # Extract work item ID and link back to bead
  AZ_WORK_ITEM_ID=$(jq -r '.id' /tmp/az_work_item_$BEAD_ID.json)
  
  if [[ -n "$AZ_WORK_ITEM_ID" && "$AZ_WORK_ITEM_ID" != "null" ]]; then
    bd update "$BEAD_ID" --metadata "az_work_item=$AZ_WORK_ITEM_ID"
    echo "✅ Created Azure DevOps work item #$AZ_WORK_ITEM_ID for bead $BEAD_ID"
  else
    echo "⚠️  Failed to create Azure DevOps work item for bead $BEAD_ID"
  fi
  
  # Clean up temp file
  rm -f /tmp/az_work_item_$BEAD_ID.json
else
  echo "Skipping Azure DevOps work item creation for bead $BEAD_ID (type: $BEAD_TYPE, priority: $BEAD_PRIORITY)"
fi
```

**Setup:**

```bash
chmod +x .beads-hooks/post-create.sh

# Test it
bd create "Test bug" -t bug -p 0
# Should automatically create AZ work item
```

---

## Status Synchronization

### Sync Bead Status to Azure DevOps

```bash
#!/bin/bash
# sync-bead-to-azure.sh
# Synchronizes bead status changes to Azure DevOps

BEAD_ID="$1"

if [[ -z "$BEAD_ID" ]]; then
  echo "Usage: $0 <bead-id>"
  exit 1
fi

# Get bead details
BEAD_JSON=$(bd show "$BEAD_ID" --json)
BEAD_STATUS=$(echo "$BEAD_JSON" | jq -r '.status')
AZ_WORK_ITEM=$(echo "$BEAD_JSON" | jq -r '.metadata.az_work_item // empty')

if [[ -z "$AZ_WORK_ITEM" ]]; then
  echo "⚠️  Bead $BEAD_ID has no linked Azure DevOps work item"
  exit 0
fi

# Map bead status to Azure DevOps state
case "$BEAD_STATUS" in
  "open")
    AZ_STATE="New"
    ;;
  "in_progress")
    AZ_STATE="Active"
    ;;
  "closed")
    AZ_STATE="Closed"
    ;;
  "blocked")
    AZ_STATE="Active"
    # Add "Blocked" tag
    az boards work-item update "$AZ_WORK_ITEM" \
      --state "$AZ_STATE" \
      --discussion "Status synced from Beads: $BEAD_STATUS"
    
    # Add blocked tag
    EXISTING_TAGS=$(az boards work-item show --id "$AZ_WORK_ITEM" --query "fields.['System.Tags']" -o tsv)
    if [[ ! "$EXISTING_TAGS" =~ "Blocked" ]]; then
      NEW_TAGS="${EXISTING_TAGS}; Blocked"
      az boards work-item update "$AZ_WORK_ITEM" --fields "System.Tags=$NEW_TAGS"
    fi
    
    echo "✅ Updated Azure DevOps work item #$AZ_WORK_ITEM to $AZ_STATE (Blocked)"
    exit 0
    ;;
  *)
    echo "⚠️  Unknown bead status: $BEAD_STATUS"
    exit 0
    ;;
esac

# Update Azure DevOps work item
az boards work-item update "$AZ_WORK_ITEM" \
  --state "$AZ_STATE" \
  --discussion "Status synced from Beads: $BEAD_STATUS"

echo "✅ Updated Azure DevOps work item #$AZ_WORK_ITEM to $AZ_STATE"
```

**Usage:**

```bash
chmod +x sync-bead-to-azure.sh

# Update bead status
bd update bd-abc --status in_progress

# Sync to Azure DevOps
./sync-bead-to-azure.sh bd-abc
```

### Bidirectional Sync Script

```bash
#!/bin/bash
# bidirectional-sync.sh
# Syncs beads and Azure DevOps work items in both directions

AZ_PROJECT="MyProject"

sync_beads_to_azure() {
  echo "🔄 Syncing beads to Azure DevOps..."
  
  # Get all open beads with Azure DevOps work items
  bd list --status open --json | jq -r '.[] | select(.metadata.az_work_item != null) | .id' | while read -r BEAD_ID; do
    ./sync-bead-to-azure.sh "$BEAD_ID"
  done
}

sync_azure_to_beads() {
  echo "🔄 Syncing Azure DevOps to beads..."
  
  # Get all active work items assigned to me
  az boards query --wiql "SELECT [System.Id], [System.Title], [System.State] FROM WorkItems WHERE [System.AssignedTo] = @Me AND [System.State] <> 'Closed'" --output json > /tmp/my_work_items.json
  
  # Process each work item
  jq -c '.[]' /tmp/my_work_items.json | while read -r WORK_ITEM; do
    WORK_ITEM_ID=$(echo "$WORK_ITEM" | jq -r '.fields["System.Id"]')
    WORK_ITEM_TITLE=$(echo "$WORK_ITEM" | jq -r '.fields["System.Title"]')
    WORK_ITEM_STATE=$(echo "$WORK_ITEM" | jq -r '.fields["System.State"]')
    
    # Check if bead exists for this work item
    EXISTING_BEAD=$(bd list --json | jq -r ".[] | select(.metadata.az_work_item == \"$WORK_ITEM_ID\") | .id")
    
    if [[ -z "$EXISTING_BEAD" ]]; then
      # Create bead if it doesn't exist
      echo "Creating bead for Azure DevOps work item #$WORK_ITEM_ID..."
      
      bd create "$WORK_ITEM_TITLE (AZ-$WORK_ITEM_ID)" -p 1 --json > /tmp/new_bead.json
      NEW_BEAD_ID=$(jq -r '.id' /tmp/new_bead.json)
      
      # Link to Azure DevOps
      bd update "$NEW_BEAD_ID" --metadata "az_work_item=$WORK_ITEM_ID"
      
      echo "✅ Created bead $NEW_BEAD_ID for Azure DevOps work item #$WORK_ITEM_ID"
    else
      echo "Bead $EXISTING_BEAD already exists for Azure DevOps work item #$WORK_ITEM_ID"
    fi
  done
  
  rm -f /tmp/my_work_items.json /tmp/new_bead.json
}

# Run both syncs
sync_beads_to_azure
sync_azure_to_beads

echo "✅ Bidirectional sync complete"
```

**Usage:**

```bash
chmod +x bidirectional-sync.sh

# Run sync (typically in a cron job or at start of day)
./bidirectional-sync.sh
```

---

## Agent Session Scripts

### Agent Startup Script

```bash
#!/bin/bash
# agent-startup.sh
# Run this at the start of an agent session

echo "🤖 Agent session starting..."

# 1. Pull latest changes
echo "📥 Pulling latest changes..."
git pull --rebase

# 2. Import any updated beads
echo "📦 Importing beads..."
bd sync

# 3. Show ready work
echo ""
echo "✅ Ready work (no blockers):"
bd ready --json | jq -r '.[] | "  - \(.id): \(.title) (P\(.priority))"'

# 4. Show my active work
echo ""
echo "🚧 My active work:"
bd list --status in_progress --json | jq -r '.[] | "  - \(.id): \(.title)"'

# 5. Sync with Azure DevOps
echo ""
echo "🔄 Syncing with Azure DevOps..."
./bidirectional-sync.sh

echo ""
echo "✅ Agent session ready. Use 'bd show <id>' to view task details."
```

### Agent Shutdown Script

```bash
#!/bin/bash
# agent-shutdown.sh
# Run this at the end of an agent session

echo "🤖 Agent session ending..."

# 1. Show work completed this session
echo ""
echo "✅ Work completed:"
bd list --status closed --json | jq -r 'sort_by(.closed) | reverse | .[:5] | .[] | "  - \(.id): \(.title)"'

# 2. Show open work in progress
echo ""
echo "🚧 Work still in progress:"
bd list --status in_progress --json | jq -r '.[] | "  - \(.id): \(.title)"'

# 3. Sync beads to git
echo ""
echo "💾 Syncing beads to git..."
bd sync

# 4. Push changes
echo ""
echo "📤 Pushing changes to remote..."
git push

# 5. Sync status to Azure DevOps
echo ""
echo "🔄 Syncing to Azure DevOps..."
./bidirectional-sync.sh

echo ""
echo "✅ Agent session ended. All work synced."
```

### Agent Task Execution Template

```bash
#!/bin/bash
# execute-bead-task.sh <bead-id>
# Complete workflow for executing a single bead task

BEAD_ID="$1"

if [[ -z "$BEAD_ID" ]]; then
  echo "Usage: $0 <bead-id>"
  exit 1
fi

# 1. Show task details
echo "📋 Task details:"
bd show "$BEAD_ID" --json | jq '{id, title, description, priority, blockers}'

# 2. Update status to in_progress
echo ""
echo "🚧 Starting work on $BEAD_ID..."
bd update "$BEAD_ID" --status in_progress
./sync-bead-to-azure.sh "$BEAD_ID"

# 3. [Agent performs work here]
echo ""
echo "🔨 Execute your task implementation here..."
echo "   (This script is a template - add your work logic)"

# 4. Ask for completion confirmation
echo ""
read -p "❓ Is the task complete? (y/n): " COMPLETE

if [[ "$COMPLETE" == "y" || "$COMPLETE" == "Y" ]]; then
  # 5. Close the bead
  read -p "📝 Completion reason: " REASON
  
  bd close "$BEAD_ID" --reason "$REASON"
  ./sync-bead-to-azure.sh "$BEAD_ID"
  
  echo "✅ Task $BEAD_ID completed and synced."
else
  echo "⚠️  Task $BEAD_ID marked in progress but not completed."
fi

# 6. Show next ready task
echo ""
echo "📋 Next ready tasks:"
bd ready --json | jq -r '.[:3] | .[] | "  - \(.id): \(.title) (P\(.priority))"'
```

---

## Multi-Agent Coordination

### Example: Epic Breakdown by Multiple Agents

```bash
#!/bin/bash
# epic-breakdown.sh
# Example: Lead agent breaks down epic, worker agents execute

# === LEAD AGENT ===
echo "👔 Lead Agent: Breaking down epic..."

# Create epic
bd create "Implement User Authentication System" -t epic -p 0 --json > /tmp/epic.json
EPIC_ID=$(jq -r '.id' /tmp/epic.json)

echo "Created epic: $EPIC_ID"

# Create Azure DevOps epic
az boards work-item create \
  --type Epic \
  --title "Implement User Authentication System" \
  --description "Bead: $EPIC_ID" \
  --project "MyProject" \
  --output json > /tmp/az_epic.json

AZ_EPIC_ID=$(jq -r '.id' /tmp/az_epic.json)
bd update "$EPIC_ID" --metadata "az_work_item=$AZ_EPIC_ID"

# Break down into sub-tasks
bd create "Implement login API endpoint" -p "$EPIC_ID" -t task --json > /tmp/task1.json
bd create "Implement JWT token generation" -p "$EPIC_ID" -t task --json > /tmp/task2.json
bd create "Implement password reset flow" -p "$EPIC_ID" -t task --json > /tmp/task3.json
bd create "Write integration tests" -p "$EPIC_ID" -t task --json > /tmp/task4.json

TASK1=$(jq -r '.id' /tmp/task1.json)
TASK2=$(jq -r '.id' /tmp/task2.json)
TASK3=$(jq -r '.id' /tmp/task3.json)
TASK4=$(jq -r '.id' /tmp/task4.json)

# Add dependencies (tests block on all others)
bd dep add "$TASK4" "$TASK1" --type blocks
bd dep add "$TASK4" "$TASK2" --type blocks
bd dep add "$TASK4" "$TASK3" --type blocks

# Sync to git
bd sync
git push

echo "✅ Epic breakdown complete. Tasks: $TASK1, $TASK2, $TASK3, $TASK4"

# === WORKER AGENT 1 ===
echo ""
echo "🤖 Worker Agent 1: Checking for ready work..."
git pull
bd ready --json | jq -r '.[] | select(.title | contains("login API")) | .id' | head -1 > /tmp/my_task.txt
MY_TASK=$(cat /tmp/my_task.txt)

if [[ -n "$MY_TASK" ]]; then
  echo "Starting work on $MY_TASK..."
  bd update "$MY_TASK" --status in_progress
  
  # [Do work...]
  sleep 2  # Simulate work
  
  bd close "$MY_TASK" --reason "Completed"
  bd sync
  git push
  
  echo "✅ Completed $MY_TASK"
fi

# === WORKER AGENT 2 ===
echo ""
echo "🤖 Worker Agent 2: Checking for ready work..."
git pull
bd ready --json | jq -r '.[] | select(.title | contains("JWT token")) | .id' | head -1 > /tmp/my_task2.txt
MY_TASK2=$(cat /tmp/my_task2.txt)

if [[ -n "$MY_TASK2" ]]; then
  echo "Starting work on $MY_TASK2..."
  bd update "$MY_TASK2" --status in_progress
  
  # [Do work...]
  sleep 2  # Simulate work
  
  bd close "$MY_TASK2" --reason "Completed"
  bd sync
  git push
  
  echo "✅ Completed $MY_TASK2"
fi

# === WORKER AGENT 3 ===
echo ""
echo "🤖 Worker Agent 3: Checking for ready work..."
git pull
bd ready --json | jq -r '.[] | select(.title | contains("password reset")) | .id' | head -1 > /tmp/my_task3.txt
MY_TASK3=$(cat /tmp/my_task3.txt)

if [[ -n "$MY_TASK3" ]]; then
  echo "Starting work on $MY_TASK3..."
  bd update "$MY_TASK3" --status in_progress
  
  # [Do work...]
  sleep 2  # Simulate work
  
  bd close "$MY_TASK3" --reason "Completed"
  bd sync
  git push
  
  echo "✅ Completed $MY_TASK3"
fi

# === TESTING AGENT ===
echo ""
echo "🧪 Testing Agent: Checking if tests are unblocked..."
git pull
bd ready --json | jq -r '.[] | select(.title | contains("integration tests")) | .id' | head -1 > /tmp/test_task.txt
TEST_TASK=$(cat /tmp/test_task.txt)

if [[ -n "$TEST_TASK" ]]; then
  echo "All dependencies complete! Starting integration tests..."
  bd update "$TEST_TASK" --status in_progress
  
  # [Run tests...]
  sleep 3  # Simulate tests
  
  bd close "$TEST_TASK" --reason "All tests passing"
  bd sync
  git push
  
  echo "✅ Epic $EPIC_ID complete!"
else
  echo "⏳ Tests still blocked. Waiting for other tasks to complete."
fi

# Cleanup
rm -f /tmp/*.json /tmp/*.txt
```

---

## PowerShell Examples

### PowerShell: Auto-create Work Items

```powershell
# auto-create-workitem.ps1
param(
    [Parameter(Mandatory=$true)]
    [string]$BeadId,
    
    [Parameter(Mandatory=$true)]
    [string]$BeadTitle,
    
    [Parameter(Mandatory=$true)]
    [string]$BeadType,
    
    [Parameter(Mandatory=$true)]
    [int]$BeadPriority
)

$AzProject = "MyProject"
$AzAssignedTo = "me@example.com"

# Only sync bugs, epics, and high-priority tasks
$shouldSync = $false

if ($BeadType -eq "epic" -or $BeadType -eq "bug") {
    $shouldSync = $true
} elseif ($BeadType -eq "task" -and ($BeadPriority -eq 0 -or $BeadPriority -eq 1)) {
    $shouldSync = $true
}

if (-not $shouldSync) {
    Write-Host "Skipping Azure DevOps work item creation for bead $BeadId"
    exit 0
}

Write-Host "Creating Azure DevOps work item for bead $BeadId..."

# Map bead priority to Azure DevOps priority
$azPriority = switch ($BeadPriority) {
    0 { 1 }  # Critical
    1 { 2 }  # High
    2 { 3 }  # Medium
    default { 4 }  # Low
}

# Map bead type to Azure DevOps work item type
$azType = switch ($BeadType) {
    "bug" { "Bug" }
    "epic" { "Epic" }
    default { "Task" }
}

# Create work item
$description = "Linked bead: $BeadId`n`nAutomatically created from Beads."

$workItem = az boards work-item create `
    --type $azType `
    --title $BeadTitle `
    --description $description `
    --priority $azPriority `
    --assigned-to $AzAssignedTo `
    --project $AzProject `
    --output json | ConvertFrom-Json

if ($workItem) {
    $azWorkItemId = $workItem.id
    
    # Link back to bead
    bd update $BeadId --metadata "az_work_item=$azWorkItemId"
    
    Write-Host "✅ Created Azure DevOps work item #$azWorkItemId for bead $BeadId" -ForegroundColor Green
} else {
    Write-Host "⚠️ Failed to create Azure DevOps work item for bead $BeadId" -ForegroundColor Yellow
}
```

### PowerShell: Agent Session Management

```powershell
# agent-session.ps1

function Start-AgentSession {
    Write-Host "🤖 Agent session starting..." -ForegroundColor Cyan
    
    # Pull latest changes
    Write-Host "📥 Pulling latest changes..."
    git pull --rebase
    
    # Sync beads
    Write-Host "📦 Importing beads..."
    bd sync
    
    # Show ready work
    Write-Host ""
    Write-Host "✅ Ready work (no blockers):" -ForegroundColor Green
    
    $readyWork = bd ready --json | ConvertFrom-Json
    foreach ($item in $readyWork) {
        Write-Host "  - $($item.id): $($item.title) (P$($item.priority))"
    }
    
    # Show active work
    Write-Host ""
    Write-Host "🚧 My active work:" -ForegroundColor Yellow
    
    $activeWork = bd list --status in_progress --json | ConvertFrom-Json
    foreach ($item in $activeWork) {
        Write-Host "  - $($item.id): $($item.title)"
    }
    
    Write-Host ""
    Write-Host "✅ Agent session ready." -ForegroundColor Green
}

function Stop-AgentSession {
    Write-Host "🤖 Agent session ending..." -ForegroundColor Cyan
    
    # Show completed work
    Write-Host ""
    Write-Host "✅ Work completed:" -ForegroundColor Green
    
    $closedWork = bd list --status closed --json | ConvertFrom-Json | Sort-Object -Property closed -Descending | Select-Object -First 5
    foreach ($item in $closedWork) {
        Write-Host "  - $($item.id): $($item.title)"
    }
    
    # Show in-progress work
    Write-Host ""
    Write-Host "🚧 Work still in progress:" -ForegroundColor Yellow
    
    $inProgressWork = bd list --status in_progress --json | ConvertFrom-Json
    foreach ($item in $inProgressWork) {
        Write-Host "  - $($item.id): $($item.title)"
    }
    
    # Sync to git
    Write-Host ""
    Write-Host "💾 Syncing beads to git..."
    bd sync
    
    # Push changes
    Write-Host ""
    Write-Host "📤 Pushing changes to remote..."
    git push
    
    Write-Host ""
    Write-Host "✅ Agent session ended. All work synced." -ForegroundColor Green
}

# Export functions
Export-ModuleMember -Function Start-AgentSession, Stop-AgentSession
```

**Usage:**

```powershell
# Import module
Import-Module .\agent-session.ps1

# Start session
Start-AgentSession

# ... do work ...

# End session
Stop-AgentSession
```

---

## Summary

These automation examples provide:

1. **Automatic work item creation** - Beads → Azure DevOps sync
2. **Status synchronization** - Bidirectional status updates
3. **Agent session management** - Startup/shutdown workflows
4. **Multi-agent coordination** - Epic breakdown and parallel execution
5. **Cross-platform support** - Bash and PowerShell examples

For more information:

- [Beads Integration Guide](../processes/beads-integration.md)
- [Azure DevOps Workflow](../processes/azure-devops-workflow.md)
- Beads documentation: https://github.com/steveyegge/beads
