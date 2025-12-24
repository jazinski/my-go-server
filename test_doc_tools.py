#!/usr/bin/env python3
"""
Quick test script to verify documentation tools work via MCP
"""

import json
import subprocess
import sys


def send_mcp_request(request):
    """Send JSON-RPC request to MCP server via stdin"""
    proc = subprocess.Popen(
        ["./my-go-server"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )

    request_str = json.dumps(request) + "\n"
    stdout, stderr = proc.communicate(input=request_str, timeout=5)

    if stderr:
        print(f"STDERR: {stderr}", file=sys.stderr)

    return json.loads(stdout) if stdout else None


def test_list_tools():
    """Test that documentation tools are registered"""
    print("🧪 Testing: List available tools...")

    request = {"jsonrpc": "2.0", "id": 1, "method": "tools/list"}

    response = send_mcp_request(request)

    if not response:
        print("❌ FAIL: No response from server")
        return False

    tools = response.get("result", {}).get("tools", [])
    tool_names = [t["name"] for t in tools]

    print(f"📋 Found {len(tools)} tools: {tool_names}")

    required_tools = ["list_documentation", "load_documentation"]
    for tool in required_tools:
        if tool in tool_names:
            print(f"✅ PASS: {tool} is registered")
        else:
            print(f"❌ FAIL: {tool} not found")
            return False

    return True


def test_list_documentation():
    """Test listing available documentation"""
    print("\n🧪 Testing: List documentation...")

    request = {
        "jsonrpc": "2.0",
        "id": 2,
        "method": "tools/call",
        "params": {"name": "list_documentation", "arguments": {}},
    }

    response = send_mcp_request(request)

    if not response:
        print("❌ FAIL: No response")
        return False

    result = response.get("result", {})
    content = result.get("content", [])

    if content and len(content) > 0:
        text = content[0].get("text", "")
        print(f"📚 Documentation list:\n{text[:500]}...")

        # Check for expected categories
        expected = ["Coding Standards", "Processes", "Architecture"]
        for category in expected:
            if category in text:
                print(f"✅ PASS: Found {category} category")
            else:
                print(f"❌ FAIL: Missing {category} category")
                return False

        return True
    else:
        print("❌ FAIL: No content returned")
        return False


def test_load_documentation():
    """Test loading specific documentation"""
    print("\n🧪 Testing: Load specific documentation...")

    request = {
        "jsonrpc": "2.0",
        "id": 3,
        "method": "tools/call",
        "params": {
            "name": "load_documentation",
            "arguments": {"path": "processes/beads-integration.md"},
        },
    }

    response = send_mcp_request(request)

    if not response:
        print("❌ FAIL: No response")
        return False

    result = response.get("result", {})
    content = result.get("content", [])

    if content and len(content) > 0:
        text = content[0].get("text", "")
        print(f"📄 Loaded document ({len(text)} chars):\n{text[:200]}...")
        print("✅ PASS: Documentation loaded successfully")
        return True
    else:
        print("❌ FAIL: No content returned")
        return False


def main():
    print("🚀 Testing MCP Documentation Tools\n")

    tests = [test_list_tools, test_list_documentation, test_load_documentation]

    passed = 0
    failed = 0

    for test in tests:
        try:
            if test():
                passed += 1
            else:
                failed += 1
        except Exception as e:
            print(f"❌ FAIL: Exception: {e}")
            failed += 1

    print(f"\n{'=' * 50}")
    print(f"📊 Results: {passed} passed, {failed} failed")
    print(f"{'=' * 50}")

    return 0 if failed == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
