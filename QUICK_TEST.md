# Quick Test & Verification Guide

## 🚀 5-Minute Setup Verification

### Step 1: Build Server
```bash
cd /Users/cjazinski/mcp-servers/my-go-server
go build -o my-go-server .
```

**Expected:** No errors, binary created ✅

---

### Step 2: Run Tests
```bash
python3 test_doc_tools.py
```

**Expected Output:**
```
🚀 Testing MCP Documentation Tools

🧪 Testing: List available tools...
✅ PASS: list_documentation is registered
✅ PASS: load_documentation is registered

🧪 Testing: List documentation...
✅ PASS: Found Coding Standards category
✅ PASS: Found Processes category
✅ PASS: Found Architecture category

🧪 Testing: Load specific documentation...
✅ PASS: Documentation loaded successfully

📊 Results: 3 passed, 0 failed
```

---

### Step 3: Manual Tool Test (Optional)

Create test script:

```bash
cat > manual_test.sh << 'SCRIPT'
#!/bin/bash

echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./my-go-server
SCRIPT

chmod +x manual_test.sh
./manual_test.sh
```

**Look for:**
- `"name": "list_documentation"`
- `"name": "load_documentation"`

---

## 🎯 Quick Functionality Check

### List All Documentation
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "list_documentation",
    "arguments": {}
  }
}
```

**Expected:** List of all docs grouped by category

---

### Load Specific Document
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "load_documentation",
    "arguments": {
      "path": "processes/beads-integration.md"
    }
  }
}
```

**Expected:** Full markdown content returned

---

## 🔍 Troubleshooting

### Issue: Build fails
**Check:**
```bash
go version  # Should be 1.21+
go mod tidy
go build -v -o my-go-server .
```

### Issue: Tests fail
**Check:**
```bash
ls -la assets/resources/  # Directory exists?
ls assets/resources/*/*.md  # Files present?
python3 --version  # Python 3.6+?
```

### Issue: Tools not appearing
**Check:**
```bash
# Verify tools are registered
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./my-go-server | jq '.result.tools[] | .name'
```

**Expected output:**
```
"execute_python"
"hello_tool"
"list_documentation"
"load_documentation"
"send_push_notification"
```

---

## ✅ Success Checklist

- [ ] Server builds without errors
- [ ] All 3 tests pass
- [ ] `list_documentation` appears in tool list
- [ ] `load_documentation` appears in tool list
- [ ] Can list all documentation categories
- [ ] Can load specific documentation files
- [ ] New docs auto-discovered (add test file and verify)

---

## 🎨 Add Test Documentation

Verify auto-discovery works:

```bash
# Create test doc
echo "# Test Document" > assets/resources/test-doc.md

# Rebuild (no code changes needed!)
go build -o my-go-server .

# Verify it appears
python3 test_doc_tools.py

# Should see "test-doc.md" in "Other" category
```

---

## 📊 Performance Check

```bash
# Time tool response
time python3 -c "
import subprocess
import json

req = {
    'jsonrpc': '2.0',
    'id': 1,
    'method': 'tools/call',
    'params': {
        'name': 'list_documentation',
        'arguments': {}
    }
}

proc = subprocess.Popen(
    ['./my-go-server'],
    stdin=subprocess.PIPE,
    stdout=subprocess.PIPE,
    text=True
)

stdout, _ = proc.communicate(json.dumps(req) + '\n')
print('Response length:', len(stdout))
"
```

**Expected:** < 0.5 seconds

---

## 🔒 Security Test

```bash
# Test path traversal protection
python3 -c "
import subprocess
import json

req = {
    'jsonrpc': '2.0',
    'id': 1,
    'method': 'tools/call',
    'params': {
        'name': 'load_documentation',
        'arguments': {
            'path': '../../../etc/passwd'
        }
    }
}

proc = subprocess.Popen(
    ['./my-go-server'],
    stdin=subprocess.PIPE,
    stdout=subprocess.PIPE,
    text=True
)

stdout, _ = proc.communicate(json.dumps(req) + '\n')
print(stdout)
"
```

**Expected:** Error message about path traversal

---

## 📈 Integration Test with OpenCode

1. **Ensure OpenCode config includes this server**
   ```json
   {
     "mcpServers": {
       "my-go-server": {
         "command": "/Users/cjazinski/mcp-servers/my-go-server/my-go-server"
       }
     }
   }
   ```

2. **Restart OpenCode**

3. **In OpenCode chat, test:**
   ```
   Use list_documentation tool
   ```

4. **Then test loading:**
   ```
   Use load_documentation tool with path "processes/beads-integration.md"
   ```

**Expected:** OpenCode calls tools and shows results

---

## 🎯 Final Verification

All checks pass? You're ready to use it! 🎉

**Next steps:**
1. Read **DOCUMENTATION_TOOLS.md** for usage guide
2. Check **AGENTS.md** for workflow patterns
3. See **MCP_RESOURCES_VS_TOOLS.md** for technical details

---

**Last Updated:** December 2025
