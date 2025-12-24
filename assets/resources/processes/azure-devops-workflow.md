# Azure DevOps Workflow & Best Practices

## 🎯 Overview

This document defines our team's Azure DevOps workflow for managing work items,
sprints, pull requests, and releases. Following these practices ensures
consistent project tracking and smooth team collaboration.

## 📋 Work Item Types & Hierarchy

### Work Item Hierarchy

```
Epic
├── Feature
│   ├── User Story
│   │   ├── Task
│   │   └── Bug
│   └── Bug
└── Feature
```

### Work Item Types Explained

**Epic:**

- Large initiatives spanning multiple sprints (1-6 months)
- Examples: "Implement Single Sign-On", "Migrate to .NET 8", "Redesign User
  Dashboard"
- Must have: Title, Description, Business Value, Acceptance Criteria
- Assigned to: Product Owner or Team Lead

**Feature:**

- Substantial functionality within an Epic (2-4 weeks)
- Examples: "OAuth2 Integration", "User Profile Management", "Payment Gateway"
- Must have: Title, Description, Acceptance Criteria, Epic link
- Assigned to: Development Team Lead

**User Story:**

- User-facing functionality that delivers value (1-5 days)
- Format: "As a [user type], I want [goal] so that [benefit]"
- Examples: "As a user, I want to reset my password so that I can regain access
  to my account"
- Must have: Title, Description, Acceptance Criteria, Story Points, Feature link
- Assigned to: Developer(s)

**Task:**

- Technical work item under a User Story (2-8 hours)
- Examples: "Create database migration", "Implement API endpoint", "Write unit
  tests"
- Must have: Title, Description, Original Estimate (hours), Parent link
- **MUST be assigned to a specific team member**

**Bug:**

- Defect in existing functionality
- Must have: Title, Description, Repro Steps, Priority, Severity
- Can be child of User Story or standalone
- **MUST be assigned to a specific team member**

## ✅ Work Item Requirements (Mandatory)

### Before Starting Any Work

1. **Work Item Must Exist** - No coding without a work item
2. **Work Item Must Be Assigned** - All Tasks and Bugs MUST have an assignee
3. **Work Item Must Be in Correct State** - Move to "Active" when starting work
4. **Work Item Must Have Description** - Clear explanation of what needs to be
   done
5. **Work Item Must Be in Current Sprint** - Don't work on future sprint items

### Task Requirements

✅ **MUST Have:**

- Clear, specific title (e.g., "Create User API endpoint", not "Work on user
  feature")
- Description explaining what needs to be done
- Original Estimate in hours
- Assigned to specific developer
- Linked to parent User Story
- Tags if applicable (e.g., "frontend", "backend", "database")

❌ **DON'T:**

- Create tasks without assignments
- Start work without moving task to "Active"
- Leave tasks in "Active" when done (move to "Closed")

### Bug Requirements

✅ **MUST Have:**

- Clear, reproducible title (e.g., "Login fails with special characters in
  password")
- Description with:
  - **Repro Steps**: Step-by-step instructions to reproduce
  - **Expected Result**: What should happen
  - **Actual Result**: What actually happens
  - **Environment**: Browser, OS, app version
- Priority (1=Critical, 2=High, 3=Medium, 4=Low)
- Severity (1=Critical, 2=High, 3=Medium, 4=Low)
- Assigned to specific developer
- Area Path and Iteration set correctly

## 🔄 Work Item States & Workflow

### User Story States

```
New → Active → Resolved → Closed
        ↓
     Removed (cancelled)
```

**New:**

- Just created, not yet started
- In backlog or future sprint

**Active:**

- Currently being worked on
- Move here when you start development
- Must have assigned developer

**Resolved:**

- Development complete, ready for testing
- Code reviewed and merged
- Waiting for QA/acceptance

**Closed:**

- Testing complete, accepted by Product Owner
- Deployed to production (or ready for deployment)

**Removed:**

- Cancelled, no longer needed
- Add comment explaining why

### Task States

```
To Do → In Progress → Done
```

**To Do:**

- Not started yet
- In current sprint but not active

**In Progress:**

- Currently working on it
- **Move here as soon as you start**
- Only one task should be "In Progress" per person at a time

**Done:**

- Work complete
- Code committed and pushed
- Tests passing
- **Move here immediately when finished**

### Bug States

```
New → Active → Resolved → Closed
  ↓
Cannot Reproduce / By Design
```

**New:** Bug reported, not yet investigated **Active:** Developer is working on
fix **Resolved:** Fix implemented, ready for verification **Closed:** Fix
verified and deployed **Cannot Reproduce:** Unable to replicate issue **By
Design:** Not a bug, working as intended

## 📅 Sprint Planning & Execution

### Sprint Cadence

- **Sprint Duration:** 2 weeks (10 working days)
- **Planning Meeting:** First day of sprint (2 hours)
- **Daily Standup:** Every day (15 minutes)
- **Sprint Review:** Last day of sprint (1 hour)
- **Sprint Retrospective:** Last day of sprint (1 hour)

### Sprint Planning Process

**Before Planning:**

1. Product Owner prioritizes backlog
2. Team reviews upcoming stories
3. Stories have acceptance criteria defined

**During Planning:**

1. Team commits to User Stories for sprint
2. Break down stories into Tasks
3. Assign story points (Fibonacci: 1, 2, 3, 5, 8, 13)
4. **Assign all Tasks to specific developers**
5. Set sprint goal

**Sprint Capacity:**

- Calculate team capacity (hours available)
- Account for: meetings, PTO, support work
- Don't over-commit (use 70-80% of capacity)

### Daily Standup Format

Each team member answers:

1. **What did I complete yesterday?** (Refer to work item IDs)
2. **What will I work on today?** (Refer to work item IDs)
3. **Any blockers or help needed?**

**Rules:**

- Keep it to 15 minutes
- Focus on work items, not detailed technical discussions
- Update work item states before standup
- Take detailed discussions offline

## 🔀 Pull Request Workflow

### Creating a Pull Request

**Prerequisites:**

1. ✅ Work item exists and is assigned to you
2. ✅ Work item is in "Active" or "In Progress" state
3. ✅ Code is complete and tested locally
4. ✅ All tests pass
5. ✅ Branch follows naming convention: `feature/AB#1234-short-description` or
   `bugfix/AB#5678-short-description`

**PR Requirements:**

**Title Format:**

```
AB#1234: Add user authentication endpoint
```

(AB# links PR to work item automatically)

**Description Template:**

```markdown
## Related Work Item

Resolves AB#1234

## Description

Brief description of what this PR accomplishes and why.

## Changes Made

- Added UserAuthController with Login/Logout endpoints
- Implemented JWT token generation
- Added unit tests for authentication logic
- Updated API documentation

## Testing Performed

- [ ] Unit tests pass (attach screenshot if applicable)
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] Tested on dev environment

## Screenshots (if UI changes)

[Add screenshots here]

## Checklist

- [ ] Code follows team style guidelines (.NET/React standards)
- [ ] All tests pass locally
- [ ] No unnecessary console logs or debug code
- [ ] Documentation updated (if needed)
- [ ] Database migrations included (if needed)
- [ ] Configuration changes documented
```

**PR Size Guidelines:**

- **Small:** < 200 lines changed (ideal, review in < 30 min)
- **Medium:** 200-500 lines changed (review in < 1 hour)
- **Large:** 500-1000 lines changed (needs detailed description)
- **Too Large:** > 1000 lines (consider splitting)

### Reviewing a Pull Request

**Reviewer Responsibilities:**

✅ **MUST Check:**

1. Code follows team coding standards (.NET Core, React)
2. Tests are adequate and pass
3. No security vulnerabilities (SQL injection, XSS, etc.)
4. Error handling is proper
5. No hardcoded secrets or sensitive data
6. Work item is linked correctly

**Review Response Time:**

- **Normal PRs:** Within 24 hours
- **Urgent PRs:** Within 4 hours (use "urgent" label)
- **Hotfixes:** Within 1 hour

**Approval Criteria:**

- ✅ At least 1 approval required (2 for critical changes)
- ✅ All comments resolved or addressed
- ✅ All automated checks pass (builds, tests)
- ✅ No merge conflicts

**Merge Strategy:**

- Use **Squash merge** for feature branches
- Keeps main branch history clean
- Preserves PR context in single commit

### After PR is Merged

1. **Delete source branch** (done automatically if configured)
2. **Update work item state** to "Resolved"
3. **Link PR to work item** (automatic if PR title has AB#)
4. **Notify QA** if testing is required

## 🏷️ Branch Strategy

### Branch Naming Convention

**Format:** `<type>/AB#<work-item-id>-<short-description>`

**Types:**

- `feature/` - New features
- `bugfix/` - Bug fixes
- `hotfix/` - Urgent production fixes
- `release/` - Release branches

**Examples:**

```
feature/AB#1234-user-authentication
bugfix/AB#5678-login-timeout
hotfix/AB#9012-critical-security-patch
release/v2.1.0
```

### Branch Protection Rules

**main branch:**

- ✅ Require pull request reviews (minimum 1)
- ✅ Require status checks to pass
- ✅ Require branches to be up to date
- ✅ No direct commits (all changes via PR)
- ✅ Require linked work items

**develop branch (if used):**

- ✅ Require pull request reviews
- ✅ Require status checks to pass
- ✅ Require linked work items

## 🚀 Release Process

### Release Planning

1. **Create Release Work Item** (Epic or Feature)
   - Title: "Release v2.1.0"
   - Description: Scope and goals
   - Target date

2. **Create Release Branch**
   ```bash
   git checkout develop
   git pull
   git checkout -b release/v2.1.0
   ```

3. **Finalize Release**
   - Update version numbers
   - Update CHANGELOG.md
   - Run full test suite
   - Create release PR: `release/v2.1.0` → `main`

4. **Tag Release**
   ```bash
   git tag -a v2.1.0 -m "Release version 2.1.0"
   git push origin v2.1.0
   ```

5. **Update Work Items**
   - Mark all items in release as "Closed"
   - Associate with release tag

### Deployment Checklist

**Before Deployment:**

- [ ] All tests pass (unit, integration, E2E)
- [ ] Code review approved
- [ ] Security scan complete (no critical issues)
- [ ] Database migrations tested
- [ ] Configuration verified
- [ ] Rollback plan documented

**After Deployment:**

- [ ] Smoke tests pass in production
- [ ] Monitor logs for errors (first 30 minutes)
- [ ] Notify stakeholders
- [ ] Update work items to "Closed"
- [ ] Update documentation

## 🏃 Daily Workflow (Developer Checklist)

### Starting Your Day

1. **Check Azure DevOps Board**
   - Review your assigned tasks
   - Check for blockers or comments

2. **Update Work Item States**
   - Move yesterday's completed tasks to "Done"
   - Move today's task to "In Progress"

3. **Attend Daily Standup**
   - Share progress and blockers

### During Development

1. **Create Feature Branch**
   ```bash
   git checkout main
   git pull
   git checkout -b feature/AB#1234-short-description
   ```

2. **Commit Regularly**
   ```bash
   git commit -m "feat(auth): implement JWT validation [AB#1234]"
   ```

3. **Push to Remote**
   ```bash
   git push -u origin feature/AB#1234-short-description
   ```

### Completing Work

1. **Create Pull Request**
   - Use PR template
   - Link work item in title: `AB#1234: Description`

2. **Request Review**
   - Add at least 1 reviewer
   - Respond to feedback promptly

3. **After Merge**
   - Update work item to "Resolved"
   - Delete feature branch
   - Move to next task

## 📊 Queries & Reports

### Common Queries

**My Work Items:**

```
Work Item Type = Task
AND Assigned To = @Me
AND State <> Closed
```

**Sprint Work:**

```
Work Item Type IN (User Story, Task, Bug)
AND Iteration = @CurrentIteration
```

**Blocked Items:**

```
Tags CONTAINS Blocked
AND State = Active
```

**Code Review Needed:**

```
Work Item Type IN (User Story, Task)
AND State = Resolved
```

## 🚨 Common Scenarios

### Scenario: Urgent Hotfix

1. Create Bug work item (Priority 1, Severity 1)
2. Assign to yourself immediately
3. Create hotfix branch from `main`
4. Fix and test thoroughly
5. Create PR with "urgent" label
6. Get review within 1 hour
7. Merge and deploy immediately
8. Update Bug to "Closed"

### Scenario: Task is Blocked

1. Add "Blocked" tag to work item
2. Add comment explaining blocker
3. Update state to "Blocked" (if available)
4. Notify team lead in standup
5. Work on different task while blocked
6. Remove "Blocked" tag when resolved

### Scenario: Requirements Changed

1. Discuss with Product Owner
2. Update User Story description
3. Add comment documenting change
4. Re-estimate if needed
5. Notify team if sprint commitment affected

## 🤖 AI Agent Integration with Beads

### Overview

**Beads (`bd`)** is a git-backed issue tracker designed for AI coding agents. It
provides persistent, structured memory that enables agents to handle
long-horizon tasks without losing context.

**Integration Philosophy:**

- **Beads** = Agent memory (fast, structured, dependency-aware)
- **Azure DevOps** = Team visibility (human-readable UI, reports, sprints)

### When to Create Both Beads and Azure DevOps Work Items

| Scenario                 | Create Bead? | Create AZ Work Item?               |
| ------------------------ | ------------ | ---------------------------------- |
| Agent breaks down work   | ✅ Always    | ✅ For human-facing tasks          |
| Agent detects bug/issue  | ✅ Always    | ✅ Always (link to bead)           |
| Agent planning sub-tasks | ✅ Always    | ⚠️ Optional (depends on scope)     |
| Human creates work item  | ⚠️ Optional  | ✅ Always                          |
| Multi-agent coordination | ✅ Required  | ✅ Required (one per epic/feature) |

### Creating Linked Beads + Work Items

**When agent discovers work:**

```bash
# 1. Create bead
bd create "Fix validation bug in UserController" -t bug -p 0 --json

# 2. Create Azure DevOps work item
az boards work-item create \
  --type Bug \
  --title "Fix validation bug in UserController" \
  --description "Linked to bead: bd-abc123" \
  --assigned-to "developer@example.com" \
  --project "MyProject"

# 3. Link them bidirectionally
bd update bd-abc123 --metadata "az_work_item=12345"
```

**When human creates work item:**

```bash
# 1. Human creates work item in Azure DevOps UI (ID: 12345)

# 2. Agent imports to beads
bd create "Implement feature X (AZ-12345)" -p 1 --json

# 3. Link in metadata
bd update bd-xyz --metadata "az_work_item=12345"
```

### Beads Workflow for Agents

**REQUIRED workflow for all AI agents:**

```bash
# Start of session - check ready work
bd ready --json

# Create beads for work breakdown
bd create "Implement feature X" -p 1 --json
bd create "Add tests for feature X" -p 1 --json

# Update status as you work
bd update bd-abc --status in_progress --json

# Close completed work
bd close bd-abc --reason "Completed" --json

# End of session sync
bd sync
```

### Benefits

**For Agents:**

- Fast, structured queries (`bd ready --json`)
- Automatic context preservation across sessions
- Dependency tracking (know when tasks are unblocked)
- Git-backed versioning (every state is recoverable)
- Zero-conflict merges in multi-agent workflows

**For Teams:**

- Azure DevOps remains source of truth for human workflows
- Beads provides agent-optimized layer
- Bidirectional sync keeps systems in harmony
- Full audit trail in both systems

### Full Documentation

See [Beads Integration Guide](beads-integration.md) for comprehensive workflow,
patterns, automation hooks, and best practices.

## 📚 Best Practices Summary

✅ **DO:**

- Create work items before coding
- Assign all tasks to specific people
- Update work item states daily
- Link commits and PRs to work items
- Keep PRs small and focused
- Review PRs within 24 hours
- Write clear commit messages
- Keep work items up to date

❌ **DON'T:**

- Start work without a work item
- Leave tasks unassigned
- Let work items sit in "Active" for days
- Create PRs without work item links
- Merge PRs without review
- Commit directly to main/develop
- Ignore failed builds or tests
- Skip standup updates

## 🔗 Azure DevOps URLs

**Boards:** `https://dev.azure.com/{organization}/{project}/_boards`\
**Repos:** `https://dev.azure.com/{organization}/{project}/_git`\
**Pipelines:** `https://dev.azure.com/{organization}/{project}/_build`\
**Work Items:** `https://dev.azure.com/{organization}/{project}/_workitems`

---

**Last Updated:** 2025-12-23\
**Version:** 1.0.0
