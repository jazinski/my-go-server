# Code Review Workflow

You are conducting a thorough code review following our team's standards and
best practices.

## 🎯 Review Objectives

Ensure code quality, maintainability, security, and adherence to team standards
before merging.

## 📋 Review Checklist

### 1. **Code Quality & Style**

Check if the code follows our team coding standards:

- [ ] **Naming Conventions**: Variables, functions, types use proper casing
      (camelCase/PascalCase)
- [ ] **Code Formatting**: Properly formatted (`dotnet format`, `prettier`,
      `eslint` applied)
- [ ] **Code Structure**: Logical organization, proper package/module structure
- [ ] **Comments**: Complex logic is documented, public APIs have doc comments
- [ ] **Dead Code**: No commented-out code, unused imports, or unreachable code
- [ ] **Magic Numbers**: Constants used instead of hardcoded values
- [ ] **Function Length**: Functions are focused and not too long (< 50 lines
      ideal)
- [ ] **File Length**: Files are reasonably sized (< 500 lines ideal)

### 2. **Functionality & Logic**

- [ ] **Correctness**: Code does what it's supposed to do
- [ ] **Edge Cases**: Handles boundary conditions, empty inputs, null/nil values
- [ ] **Error Handling**: All errors are checked and handled appropriately
- [ ] **Context Usage**: Context passed correctly for cancellation/timeouts (Go)
- [ ] **Resource Cleanup**: defer statements, try-finally blocks used properly
- [ ] **Concurrency**: Race conditions avoided, proper synchronization if
      concurrent
- [ ] **Performance**: No obvious performance issues (N+1 queries, inefficient
      loops)

### 3. **Testing**

- [ ] **Test Coverage**: New code has appropriate unit tests
- [ ] **Test Quality**: Tests actually test the logic, not just achieve coverage
- [ ] **Test Names**: Descriptive test names that explain what's being tested
- [ ] **Edge Cases**: Tests cover error cases and boundary conditions
- [ ] **Integration Tests**: Added/updated if needed for API/database changes
- [ ] **No Flaky Tests**: Tests are deterministic and reliable
- [ ] **Tests Pass**: All tests pass locally and in CI

### 4. **Security**

- [ ] **Input Validation**: All external input is validated and sanitized
- [ ] **SQL Injection**: Parameterized queries used, no string concatenation
- [ ] **XSS Prevention**: Output encoding applied for web responses
- [ ] **Authentication**: Proper auth checks on protected endpoints
- [ ] **Authorization**: Users can only access resources they're allowed to
- [ ] **Sensitive Data**: No passwords, API keys, tokens in code or logs
- [ ] **Crypto**: Using crypto/rand, not math/rand for security purposes
- [ ] **Dependencies**: No known vulnerabilities in new dependencies

### 5. **API & Interface Design**

- [ ] **Backward Compatibility**: Breaking changes are documented and approved
- [ ] **API Consistency**: Follows existing API patterns and conventions
- [ ] **Error Responses**: Proper HTTP status codes, consistent error format
- [ ] **Versioning**: API versioning strategy followed if applicable
- [ ] **Documentation**: API changes documented (OpenAPI/Swagger/comments)

### 6. **Database & Data**

- [ ] **Schema Changes**: Migrations are reversible and tested
- [ ] **Query Performance**: Indexes exist for frequently queried columns
- [ ] **Transactions**: ACID properties maintained where needed
- [ ] **Data Validation**: Constraints and validations at database level
- [ ] **No Data Loss**: Careful with DELETE/DROP operations

### 7. **Documentation**

- [ ] **README Updated**: If setup/usage changes
- [ ] **CHANGELOG Updated**: For version-controlled projects
- [ ] **API Docs Updated**: For API changes
- [ ] **Architecture Docs**: If design/architecture changes
- [ ] **Migration Guide**: For breaking changes

### 8. **Dependencies & Configuration**

- [ ] **Dependency Justification**: New dependencies are necessary and
      maintained
- [ ] **Version Pinning**: Dependencies pinned to specific versions
- [ ] **Configuration**: New config properly documented with examples
- [ ] **Environment Variables**: Sensible defaults, documented required vars
- [ ] **Secrets Management**: Secrets not hardcoded, proper env var usage

## 🔍 Review Process

### Step 1: Understand the Context

1. Read the PR description and linked issues
2. Understand the problem being solved
3. Review the overall approach before diving into code

### Step 2: High-Level Review

1. Check the file structure and organization
2. Verify the general architecture and design patterns
3. Ensure the solution is appropriate for the problem

### Step 3: Detailed Code Review

Go through the checklist above systematically for each changed file.

### Step 4: Test the Changes

1. Pull the branch locally
2. Run the tests: `go test ./...` or equivalent
3. Test manually if it's a user-facing change
4. Verify edge cases and error scenarios

### Step 5: Provide Feedback

Use these feedback categories:

- **🚨 MUST FIX**: Blocking issues (bugs, security, broken tests)
- **⚠️ SHOULD FIX**: Important but not blocking (style, minor issues)
- **💡 SUGGESTION**: Nice to have improvements (optional)
- **❓ QUESTION**: Clarification needed
- **✅ PRAISE**: Highlight good practices

### Example Feedback:

````markdown
## 🚨 MUST FIX

**File: `service/user.go:45`**

```go
user, _ := repo.GetUser(id) // Error is ignored
```
````

Error must be handled. If user is not found, this will panic.

**Suggested Fix:**

```go
user, err := repo.GetUser(id)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}
```

---

## ⚠️ SHOULD FIX

**File: `handler/user.go:23`** Function `HandleUserRequest` is 120 lines long.
Consider extracting validation and business logic into separate functions for
better readability and testability.

---

## 💡 SUGGESTION

**File: `repository/user.go:67`** Consider using a connection pool here for
better performance under load. The current implementation creates a new
connection for each request.

---

## ❓ QUESTION

**File: `service/auth.go:34`** Why is the timeout set to 60 seconds? Is this
intentional or should it use the default 30s timeout?

---

## ✅ PRAISE

**File: `service/user_test.go`** Excellent test coverage with table-driven
tests! The edge cases are well covered. 👍

````
## 🎨 Feedback Best Practices

### DO:

- **Be specific**: Point to exact line numbers and files
- **Be constructive**: Suggest improvements, not just problems
- **Be respectful**: Assume good intent, use collaborative language
- **Explain why**: Don't just say "change this", explain the reasoning
- **Provide examples**: Show better alternatives when possible
- **Praise good work**: Acknowledge good practices and clever solutions

### DON'T:

- **Be vague**: "This looks wrong" without specifics
- **Be condescending**: "You should know this" or "Obviously..."
- **Nitpick style**: If it's auto-formatted, don't comment on style
- **Block on preferences**: Personal preferences shouldn't block merge
- **Rewrite their code**: Suggest changes, don't demand your exact solution

## ⏱️ Review Response Times

- **Normal PRs**: Review within 24 hours
- **Urgent PRs**: Review within 2 hours (marked with `urgent` label)
- **Hotfixes**: Review immediately (within 30 minutes)

## ✅ Approval Criteria

Approve the PR if:

1. ✅ All "MUST FIX" items are resolved
2. ✅ Tests pass
3. ✅ No security vulnerabilities
4. ✅ Follows team coding standards
5. ✅ Documentation is adequate

Request changes if:

1. ❌ Critical bugs exist
2. ❌ Security issues present
3. ❌ Tests are failing
4. ❌ Breaking changes without discussion
5. ❌ Code quality significantly below standards

## 🔄 After Review

**As Reviewer:**
- Respond to author's questions
- Re-review after changes
- Approve once satisfied

**As Author:**
- Address all feedback
- Ask questions if unclear
- Update PR description if scope changed
- Re-request review after fixes

## 📚 Additional Guidelines

- **Small PRs**: Easier to review, review them first
- **Large PRs**: Take more time, consider multiple review sessions
- **Learning Opportunity**: Reviews are learning for both author and reviewer
- **Team Knowledge**: Reviews spread knowledge across the team
- **Quality Gate**: Reviews maintain code quality and prevent technical debt

---

## 🎯 Final Output Format

After completing the review, provide a summary:

```markdown
## Code Review Summary

### Overall Assessment
[Brief overview of the PR and changes]

### Strengths
- [Good aspect 1]
- [Good aspect 2]

### Issues Found
- 🚨 [Critical issue count] must-fix items
- ⚠️ [Important issue count] should-fix items
- 💡 [Suggestion count] suggestions

### Recommendation
- [ ] ✅ Approve (ready to merge)
- [ ] 🔄 Request Changes (issues must be addressed)
- [ ] 💬 Comment (feedback provided, but not blocking)

### Next Steps
[What the author should do next]
````

Remember: Code reviews are collaborative, not confrontational. The goal is to
improve code quality and share knowledge, not to criticize the author.
