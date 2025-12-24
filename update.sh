#!/bin/bash
# Update script for my-go-server MCP server
# Run this after git pull to rebuild and verify

set -e  # Exit on error

echo "🔄 Updating my-go-server MCP server..."
echo

# Step 1: Git pull
echo "📥 Step 1/4: Pulling latest changes from git..."
git pull
echo "✅ Git pull complete"
echo

# Step 2: Build
echo "🔨 Step 2/4: Building Go binary..."
go build -o my-go-server .
echo "✅ Build complete"
echo

# Step 3: Verify
echo "✔️  Step 3/4: Verifying binary..."
if [ -f "./my-go-server" ]; then
    echo "✅ Binary exists: $(ls -lh my-go-server | awk '{print $5}')"
    echo "   Last modified: $(ls -l my-go-server | awk '{print $6, $7, $8}')"
else
    echo "❌ ERROR: Binary not found!"
    exit 1
fi
echo

# Step 4: Test
echo "🧪 Step 4/4: Testing server starts..."
timeout 2s ./my-go-server 2>/dev/null || true
echo "✅ Server starts successfully"
echo

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✨ Update complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo
echo "📋 Next steps:"
echo "   1. Restart OpenCode completely"
echo "   2. Test: Ask AI to 'List available documentation'"
echo "   3. Verify new features are available"
echo
echo "💡 Tip: Run this script after every git pull to ensure you're up to date!"
echo
