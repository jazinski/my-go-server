# Git Workflow & Branch Strategy

## 🎯 Overview

This document defines our team's Git workflow, branching strategy, commit
conventions, and pull request process. Following these guidelines ensures a
clean git history and smooth collaboration.

## 🌳 Branch Strategy

We use a **simplified Git Flow** model:

```
main (production)
  ├── develop (integration)
  │   ├── feature/user-authentication
  │   ├── feature/api-pagination
  │   ├── bugfix/login-timeout
  │   └── hotfix/critical-security-patch
  └── release/v1.2.0
```

### Branch Types

| Branch Type | Naming                  | Purpose                 | Base Branch | Merge To             |
| ----------- | ----------------------- | ----------------------- | ----------- | -------------------- |
| **main**    | `main`                  | Production code         | -           | -                    |
| **develop** | `develop`               | Integration branch      | `main`      | `main` (via release) |
| **feature** | `feature/<description>` | New features            | `develop`   | `develop`            |
| **bugfix**  | `bugfix/<description>`  | Bug fixes               | `develop`   | `develop`            |
| **hotfix**  | `hotfix/<description>`  | Urgent production fixes | `main`      | `main` + `develop`   |
| **release** | `release/v<version>`    | Release preparation     | `develop`   | `main` + `develop`   |

### Branch Naming Rules

✅ **Good Examples:**

```
feature/user-authentication
feature/oauth2-integration
bugfix/memory-leak-in-parser
hotfix/sql-injection-vulnerability
release/v1.2.0
```

❌ **Bad Examples:**

```
new-feature              # Not descriptive
fix                      # Too generic
johns-branch             # Personal ownership
FEATURE-123              # Use lowercase
feature_user_auth        # Use hyphens, not underscores
```

### Naming Convention:

- Use **lowercase**
- Use **hyphens** to separate words
- Be **descriptive** but concise (3-5 words max)
- **Azure DevOps**: Include work item ID: `feature/AB#1234-user-auth`
- Format: `<type>/AB#<work-item-id>-<short-description>`

## 📝 Commit Message Convention

We follow **Conventional Commits** specification:

### Format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types:

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, missing semicolons, etc.)
- **refactor**: Code refactoring (no functional changes)
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Build process, tooling, dependencies

### Examples:

```
feat(auth): add OAuth2 authentication [AB#1234]

Implement OAuth2 authentication flow with support for
Google and GitHub providers. Users can now sign in using
their existing accounts.

Related Work Item: AB#1234
```

```
fix(api): resolve timeout issue in user endpoint [AB#456]

The /api/users endpoint was timing out for requests with
large result sets. Added pagination and database query
optimization to resolve the issue.

Resolves AB#456
```

```
refactor(database): simplify connection pool management [AB#789]

Extracted connection pool logic into separate package for
better reusability and testing.

Related Work Item: AB#789
```

```
chore(deps): update .NET dependencies [AB#999]

- Update Entity Framework Core to v8.0
- Update ASP.NET Core to v8.0
- Update Newtonsoft.Json to v13.0.3

Related Work Item: AB#999
```

### Rules:

1. **First line (subject)**:
   - Max 72 characters
   - Use imperative mood: "add" not "added" or "adds"
   - No period at the end
   - Lowercase after colon

2. **Body** (optional but recommended):
   - Wrap at 72 characters
   - Explain WHAT and WHY, not HOW
   - Leave one blank line between subject and body

3. **Footer** (optional):
   - **Azure DevOps**: Reference work items: `Related Work Item: AB#123`,
     `Resolves AB#456`, `Refs AB#789`
   - Breaking changes: `BREAKING CHANGE: description`
   - **Note**: Azure DevOps auto-links work items when you use `AB#` syntax

## 🔄 Workflow Process

### 1. Starting New Work

```bash
# Update your local develop branch
git checkout develop
git pull origin develop

# Create new feature branch
git checkout -b feature/user-authentication

# Make changes and commit regularly
git add .
git commit -m "feat(auth): implement basic auth middleware [AB#1234]"

# Push to remote regularly (backup + visibility)
git push -u origin feature/AB#1234-user-authentication
```

### 2. During Development

```bash
# Keep your branch up to date with develop
git checkout develop
git pull origin develop
git checkout feature/user-authentication
git rebase develop  # or merge if you prefer

# If conflicts occur
git status
# Resolve conflicts in files
git add <resolved-files>
git rebase --continue

# Force push after rebase (your branch only!)
git push --force-with-lease origin feature/user-authentication
```

### 3. Ready for Review

```bash
# Ensure branch is up to date
git checkout develop
git pull origin develop
git checkout feature/user-authentication
git rebase develop

# Clean up commits if needed (squash work-in-progress commits)
git rebase -i develop

# Push and create pull request
git push --force-with-lease origin feature/user-authentication
# Then create PR in GitHub/GitLab
```

### 4. After PR Approval

```bash
# Merge via GitHub/GitLab UI (squash or merge commit based on preference)
# Delete remote branch via UI

# Clean up local branches
git checkout develop
git pull origin develop
git branch -d feature/user-authentication
```

## 🔍 Pull Request Process

### Creating a PR

1. **Title**: Follow commit message convention + link work item
   ```
   AB#1234: Add OAuth2 authentication (feat)
   AB#5678: Fix timeout in user endpoint (fix)
   ```
   **Note**: Including `AB#` in title auto-links PR to work item in Azure DevOps

2. **Description Template**:
   ```markdown
   ## Related Work Item

   Resolves AB#123

   ## Description

   Brief description of what this PR does and why.

   ## Changes

   - Change 1
   - Change 2
   - Change 3

   ## Testing

   - [ ] Unit tests added/updated
   - [ ] Integration tests added/updated
   - [ ] Manual testing completed
   - [ ] Documentation updated

   ## Screenshots (if applicable)

   [Add screenshots for UI changes]

   ## Checklist

   - [ ] Code follows team style guidelines
   - [ ] Self-review completed
   - [ ] Comments added for complex logic
   - [ ] No unnecessary console.log/Debug.WriteLine statements
   - [ ] Tests pass locally
   - [ ] Documentation updated
   - [ ] Work item linked (AB# in title)
   ```

3. **Size Guidelines**:
   - Keep PRs small: **< 400 lines changed** (ideal)
   - Large PRs (> 800 lines) should be split if possible
   - If unavoidable, add detailed description and comments

### Reviewing a PR

**As a Reviewer:**

1. **Check for:**
   - Code follows team standards
   - Tests are adequate
   - No obvious bugs or security issues
   - Clear variable/function names
   - Proper error handling
   - Documentation is updated

2. **Review Types:**
   - **Approve**: Ready to merge
   - **Request Changes**: Issues that must be fixed
   - **Comment**: Suggestions or questions (doesn't block merge)

3. **Response Time:**
   - Review PRs within **24 hours**
   - Urgent PRs: within **2 hours** (communicate in chat)

4. **Feedback Style:**
   ```
   # ✅ Good feedback (constructive):
   "Consider using a map here instead of a slice for O(1) lookups. 
   With the current approach, performance degrades with list size."

   # ❌ Poor feedback (unconstructive):
   "This is wrong."
   "Why did you do it this way?"
   ```

**As an Author:**

1. **Respond to all comments** (even if just "Fixed" or "Done")
2. **Push fixes in new commits** (don't force push during review)
3. **Re-request review** after addressing feedback
4. **Ask questions** if feedback is unclear

### Merge Strategies

**Squash and Merge** (Preferred for feature branches):

- Combines all commits into one
- Keeps `develop` history clean
- Use when branch has many "WIP" commits

**Merge Commit** (For release branches):

- Preserves all commits
- Shows clear merge point
- Use for `release/*` → `main` merges

**Rebase and Merge** (Avoid unless necessary):

- Creates linear history
- Can cause confusion with multiple contributors

## 🚨 Hotfix Process

For critical production issues:

```bash
# 1. Create hotfix from main
git checkout main
git pull origin main
git checkout -b hotfix/critical-sql-injection

# 2. Fix the issue
git add .
git commit -m "fix(security): patch SQL injection vulnerability"

# 3. Push and create PR to main
git push -u origin hotfix/critical-sql-injection
# Create PR: hotfix/critical-sql-injection → main

# 4. After merge to main, also merge to develop
git checkout develop
git pull origin develop
git merge main
git push origin develop

# 5. Clean up
git branch -d hotfix/critical-sql-injection
```

## 🏷️ Release Process

```bash
# 1. Create release branch from develop
git checkout develop
git pull origin develop
git checkout -b release/v1.2.0

# 2. Update version numbers, CHANGELOG, etc.
git add .
git commit -m "chore(release): prepare v1.2.0"

# 3. Create PR: release/v1.2.0 → main
# After approval and merge:

# 4. Tag the release
git checkout main
git pull origin main
git tag -a v1.2.0 -m "Release version 1.2.0"
git push origin v1.2.0

# 5. Merge back to develop
git checkout develop
git merge main
git push origin develop

# 6. Clean up
git branch -d release/v1.2.0
```

## 🛠️ Git Configuration

### Required Setup

```bash
# Set your identity
git config --global user.name "Your Name"
git config --global user.email "your.email@company.com"

# Use main as default branch name
git config --global init.defaultBranch main

# Enable colored output
git config --global color.ui auto

# Set default editor (choose one)
git config --global core.editor "code --wait"  # VS Code
git config --global core.editor "vim"          # Vim

# Pull with rebase by default (cleaner history)
git config --global pull.rebase true

# Use force-with-lease instead of force (safer)
git config --global alias.pushf 'push --force-with-lease'
```

### Useful Aliases

```bash
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.st status
git config --global alias.unstage 'reset HEAD --'
git config --global alias.last 'log -1 HEAD'
git config --global alias.visual 'log --oneline --graph --decorate --all'
```

## 🚫 Common Mistakes to Avoid

### ❌ DON'T:

1. **Commit directly to `main` or `develop`**
   ```bash
   git checkout develop  # ❌
   git commit -m "quick fix"
   git push
   ```

2. **Use generic commit messages**
   ```bash
   git commit -m "fix"          # ❌
   git commit -m "updates"      # ❌
   git commit -m "WIP"          # ❌ (not in final PR)
   ```

3. **Force push to shared branches**
   ```bash
   git push --force origin develop  # ❌ NEVER!
   git push --force origin main     # ❌ NEVER!
   ```

4. **Commit sensitive data**
   ```bash
   # ❌ NEVER commit:
   - API keys, passwords, tokens
   - .env files with secrets
   - Private certificates/keys
   - Customer data
   ```

5. **Create long-lived feature branches**
   ```bash
   # Feature branch open for 3 months  # ❌
   # Merge conflicts nightmare
   ```

### ✅ DO:

1. **Always work in feature branches**
   ```bash
   git checkout -b feature/my-feature  # ✅
   ```

2. **Write descriptive commit messages**
   ```bash
   git commit -m "feat(auth): add JWT token validation"  # ✅
   ```

3. **Keep branches up to date**
   ```bash
   git rebase develop  # ✅ Regular updates
   ```

4. **Use .gitignore properly**
   ```gitignore
   # ✅ DO ignore:
   .env
   .env.local
   *.key
   *.pem
   node_modules/
   vendor/
   ```

5. **Merge/rebase frequently**
   ```bash
   # ✅ Small, frequent PRs
   # Easier to review, fewer conflicts
   ```

## 🔐 Security Checklist

Before committing:

- [ ] No API keys, passwords, or tokens
- [ ] No customer/user data
- [ ] No private certificates or keys
- [ ] `.env` files are in `.gitignore`
- [ ] No commented-out sensitive code
- [ ] No internal URLs or endpoints (use environment variables)

**If you accidentally commit sensitive data:**

```bash
# Contact team lead immediately
# Use git filter-branch or BFG Repo-Cleaner to remove from history
# Rotate all exposed credentials
```

## 📚 Additional Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Git Flow](https://nvie.com/posts/a-successful-git-branching-model/)
- [GitHub Flow](https://guides.github.com/introduction/flow/)
- [How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/)
- [Git Best Practices](https://git-scm.com/book/en/v2)

---

**Last Updated:** 2025-01-23 **Version:** 1.0.0
