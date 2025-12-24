#!/bin/bash
# Diagnostic script for my-go-server MCP server
# Helps troubleshoot common issues

echo "🔍 my-go-server Diagnostic Tool"
echo "================================"
echo

# Check 1: Binary exists
echo "✓ Check 1: Binary exists"
if [ -f "./my-go-server" ]; then
    echo "  ✅ Binary found: $(ls -lh my-go-server | awk '{print $5}')"
    echo "     Last modified: $(ls -l my-go-server | awk '{print $6, $7, $8}')"
else
    echo "  ❌ Binary NOT found!"
    echo "     Run: go build -o my-go-server ."
    exit 1
fi
echo

# Check 2: Binary is executable
echo "✓ Check 2: Binary is executable"
if [ -x "./my-go-server" ]; then
    echo "  ✅ Binary is executable"
else
    echo "  ❌ Binary is not executable!"
    echo "     Run: chmod +x my-go-server"
    exit 1
fi
echo

# Check 3: Assets directory exists
echo "✓ Check 3: Assets directory structure"
if [ -d "./assets" ]; then
    echo "  ✅ assets/ directory exists"
    
    if [ -d "./assets/resources" ]; then
        RESOURCE_COUNT=$(find ./assets/resources -name "*.md" | wc -l | tr -d ' ')
        echo "  ✅ assets/resources/ exists ($RESOURCE_COUNT markdown files)"
    else
        echo "  ❌ assets/resources/ directory missing!"
        exit 1
    fi
    
    if [ -d "./assets/prompts" ]; then
        PROMPT_COUNT=$(find ./assets/prompts -name "*.md" | wc -l | tr -d ' ')
        echo "  ✅ assets/prompts/ exists ($PROMPT_COUNT markdown files)"
    else
        echo "  ❌ assets/prompts/ directory missing!"
        exit 1
    fi
else
    echo "  ❌ assets/ directory NOT found!"
    echo "     Make sure you're running this from the my-go-server directory"
    echo "     Current directory: $(pwd)"
    exit 1
fi
echo

# Check 4: Sample file check
echo "✓ Check 4: Sample documentation files"
SAMPLE_FILES=(
    "assets/resources/coding-standards/reactjs-style-guide.md"
    "assets/resources/processes/beads-integration.md"
    "assets/prompts/code-review.md"
)

ALL_SAMPLES_EXIST=true
for file in "${SAMPLE_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✅ $file"
    else
        echo "  ❌ $file (missing)"
        ALL_SAMPLES_EXIST=false
    fi
done

if [ "$ALL_SAMPLES_EXIST" = false ]; then
    echo
    echo "  ⚠️  Some files are missing. Did you git pull?"
    exit 1
fi
echo

# Check 5: Git status
echo "✓ Check 5: Git repository status"
if [ -d ".git" ]; then
    echo "  ✅ Git repository detected"
    
    # Check for uncommitted changes
    if git diff-index --quiet HEAD -- 2>/dev/null; then
        echo "  ✅ No uncommitted changes"
    else
        echo "  ⚠️  You have uncommitted changes"
    fi
    
    # Check if behind remote
    git fetch --dry-run 2>&1 | grep -q "From" && {
        LOCAL=$(git rev-parse @)
        REMOTE=$(git rev-parse @{u} 2>/dev/null)
        if [ "$LOCAL" = "$REMOTE" ]; then
            echo "  ✅ Up to date with remote"
        else
            echo "  ⚠️  Your local branch is behind remote"
            echo "     Run: ./update.sh"
        fi
    }
else
    echo "  ⚠️  Not a git repository (that's okay if you downloaded as zip)"
fi
echo

# Check 6: Go environment
echo "✓ Check 6: Go environment"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo "  ✅ Go installed: $GO_VERSION"
else
    echo "  ❌ Go not found in PATH"
    echo "     Install from: https://go.dev/doc/install"
fi
echo

# Summary
echo "================================"
echo "✨ Diagnostic Summary"
echo "================================"
echo
echo "Current directory: $(pwd)"
echo "Binary location:   $(pwd)/my-go-server"
echo "Assets location:   $(pwd)/assets"
echo
echo "📋 For OpenCode configuration, use:"
echo
echo '{'
echo '  "mcpServers": {'
echo '    "team-standards": {'
echo "      \"command\": \"$(pwd)/my-go-server\","
echo '      "args": []'
echo '    }'
echo '  }'
echo '}'
echo
echo "💡 Next steps:"
echo "   1. Copy the configuration above to your OpenCode config"
echo "   2. Restart OpenCode completely"
echo "   3. Test: 'List available documentation'"
echo
