# Playwright Configuration

## Overview

This MCP server uses a Playwright-enabled Docker image for Python execution,
allowing agents to take screenshots and scrape web content.

## Docker Image

- **Image**: `mcr.microsoft.com/playwright/python:v1.57.0-noble`
- **Playwright Version**: 1.57.0 (automatically pinned)
- **Python Version**: 3.12.3
- **Supported Browsers**: Chromium, Firefox, WebKit

## Key Features

1. **Automatic Version Pinning**: When you request the `playwright` module, the
   server automatically installs `playwright==1.57.0` to match the Docker
   image's browser versions.

2. **Headless Browsers**: All browsers are pre-installed and ready to use in
   headless mode.

3. **Screenshot Support**: Save screenshots to `/output` directory which maps to
   `python_output/` on the host.

## Usage Example

```python
from playwright.sync_api import sync_playwright

with sync_playwright() as p:
    browser = p.chromium.launch(headless=True)
    page = browser.new_page()
    
    # Navigate to any website
    page.goto("https://smoke-stack.app", wait_until="networkidle", timeout=30000)
    
    # Take screenshot
    page.screenshot(path="/output/screenshot.png", full_page=True)
    
    # Print results
    print(f"Title: {page.title()}")
    print(f"URL: {page.url}")
    
    browser.close()
```

## Common Patterns

### Taking Screenshots

```python
# Full page screenshot
page.screenshot(path="/output/fullpage.png", full_page=True)

# Viewport only
page.screenshot(path="/output/viewport.png")

# Specific element
element = page.locator("#my-element")
element.screenshot(path="/output/element.png")
```

### Waiting for Content

```python
# Wait for network to be idle
page.goto(url, wait_until="networkidle")

# Wait for specific element
page.wait_for_selector("#content")

# Wait for timeout
page.wait_for_timeout(2000)  # 2 seconds
```

### Interacting with Pages

```python
# Click elements
page.click("button#submit")

# Fill forms
page.fill("input[name='email']", "test@example.com")

# Get text content
content = page.text_content(".main-content")
print(content)
```

## Troubleshooting

### Issue: "Executable doesn't exist" Error

**Cause**: Playwright version mismatch between pip package and Docker image
browsers.

**Solution**: The server now automatically pins `playwright==1.57.0` when you
request it as a module. Simply use `modules="playwright"` and the correct
version will be installed.

### Issue: Timeout Errors

**Cause**: Page takes too long to load or network requests are slow.

**Solution**:

- Increase timeout: `page.goto(url, timeout=60000)` (60 seconds)
- Use different wait strategy: `wait_until="domcontentloaded"` instead of
  `"networkidle"`

### Issue: Screenshots Not Saved

**Cause**: Incorrect path - must use `/output` directory.

**Solution**: Always save to `/output/filename.png`, which maps to
`python_output/` on the host.

## Version History

- **v1.57.0** (Current): Updated from v1.49.1 to fix browser executable issues.
  Automatic version pinning implemented to prevent version mismatches.

## Best Practices

1. **Always use print statements**: Output is captured from stdout/stderr
2. **Save files to /output**: This directory persists after execution
3. **Use appropriate timeouts**: Default is 30 seconds for navigation
4. **Close browsers**: Always `browser.close()` to free resources
5. **Handle errors**: Use try/except blocks for robust scripts
6. **Request 'playwright' module**: The server handles version pinning
   automatically

## Module Request

When using the `execute_python` tool, simply specify:

```json
{
  "code": "your python code here",
  "modules": "playwright"
}
```

The server will automatically install `playwright==1.57.0` to ensure
compatibility with the Docker image browsers.

## Additional Resources

- [Playwright Python Documentation](https://playwright.dev/python/)
- [Microsoft Playwright Docker Images](https://mcr.microsoft.com/en-us/product/playwright/python/about)
